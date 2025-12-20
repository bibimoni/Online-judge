import { Test, TestingModule } from '@nestjs/testing';
import { AuthService } from './auth.service';
import { PrismaService } from '../prisma/prisma.service';
import { JwtService } from '@nestjs/jwt';
import { ConflictException, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcrypt';

describe('AuthService', () => {
  let service: AuthService;
  let prisma: PrismaService;
  let jwtService: JwtService;

  const mockPrismaService = {
    user: {
      create: jest.fn(),
      findUnique: jest.fn(),
      update: jest.fn(),
    },
  };

  const mockJwtService = {
    sign: jest.fn(),
    verify: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
        {
          provide: JwtService,
          useValue: mockJwtService,
        },
      ],
    }).compile();

    service = module.get<AuthService>(AuthService);
    prisma = module.get<PrismaService>(PrismaService);
    jwtService = module.get<JwtService>(JwtService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('register', () => {
    it('should register a new user successfully', async () => {
      const registerDto = {
        email: 'test@example.com',
        username: 'testuser',
        password: 'SecurePass123!',
      };

      mockPrismaService.user.findUnique.mockResolvedValue(null);
      mockPrismaService.user.create.mockResolvedValue({
        id: '1',
        email: registerDto.email,
        username: registerDto.username,
        passwordHash: 'hashed',
        createdAt: new Date(),
      });

      const result = await service.register(registerDto);

      const user = await prisma.user.findUnique({ where: { id: result.user_id } });

      expect(result).toHaveProperty('id');
      expect(user.email).toBe(registerDto.email);
      expect(mockPrismaService.user.create).toHaveBeenCalled();
    });

    it('should throw ConflictException if email exists', async () => {
      const registerDto = {
        email: 'existing@example.com',
        username: 'testuser',
        password: 'SecurePass123!',
      };

      mockPrismaService.user.findUnique.mockResolvedValue({
        id: '1',
        email: registerDto.email,
      });

      await expect(service.register(registerDto)).rejects.toThrow(
        ConflictException,
      );
    });

    it('should hash password before saving', async () => {
      const registerDto = {
        email: 'test@example.com',
        username: 'testuser',
        password: 'PlainTextPassword',
      };

      mockPrismaService.user.findUnique.mockResolvedValue(null);
      mockPrismaService.user.create.mockImplementation((data) => {
        expect(data.data.passwordHash).not.toBe(registerDto.password);
        return Promise.resolve({ id: '1', ...data.data });
      });

      await service.register(registerDto);
      expect(mockPrismaService.user.create).toHaveBeenCalled();
    });
  });

  describe('login', () => {
    it('should return tokens on successful login', async () => {
      const loginDto = {
        username: 'testuser',
        password: 'SecurePass123!',
      };

      const user = {
        id: '1',
        username: 'testuser',
        passwordHash: await bcrypt.hash(loginDto.password, 10),
      };

      mockPrismaService.user.findUnique.mockResolvedValue(user);
      mockJwtService.sign.mockReturnValue('mock-token');

      const result = await service.login(loginDto);

      expect(result).toHaveProperty('accessToken');
      expect(mockJwtService.sign).toHaveBeenCalledTimes(2);
    });

    it('should throw UnauthorizedException for invalid credentials', async () => {
      const loginDto = {
        username: 'testuser',
        password: 'WrongPassword',
      };

      mockPrismaService.user.findUnique.mockResolvedValue(null);

      await expect(service.login(loginDto)).rejects.toThrow(
        UnauthorizedException,
      );
    });

    it('should throw UnauthorizedException for wrong password', async () => {
      const loginDto = {
        username: 'testuser',
        password: 'WrongPassword',
      };

      const user = {
        id: '1',
        username: loginDto.username,
        passwordHash: await bcrypt.hash('CorrectPassword', 10),
      };

      mockPrismaService.user.findUnique.mockResolvedValue(user);

      await expect(service.login(loginDto)).rejects.toThrow(
        UnauthorizedException,
      );
    });
  });

  // describe('validateUser', () => {
  //   it('should return user if validation succeeds', async () => {
  //     const userId = '1';
  //     const user = {
  //       id: userId,
  //       email: 'test@example.com',
  //       username: 'testuser',
  //     };
  //
  //     mockPrismaService.user.findUnique.mockResolvedValue(user);
  //
  //     const result = await service.validateUser(userId);
  //
  //     expect(result).toEqual(user);
  //   });
  //
  //   it('should return null if user not found', async () => {
  //     mockPrismaService.user.findUnique.mockResolvedValue(null);
  //
  //     const result = await service.validateUser('999');
  //
  //     expect(result).toBeNull();
  //   });
  // });

  // describe('refreshToken', () => {
  //   it('should return new access token', async () => {
  //     const refreshToken = 'valid-refresh-token';
  //     const payload = { sub: '1', email: 'test@example.com' };
  //
  //     mockJwtService.verify.mockReturnValue(payload);
  //     mockJwtService.sign.mockReturnValue('new-access-token');
  //
  //     const result = await service.refreshToken(refreshToken);
  //
  //     expect(result).toHaveProperty('accessToken');
  //     expect(mockJwtService.verify).toHaveBeenCalledWith(refreshToken);
  //   });
  //
  //   it('should throw UnauthorizedException for invalid token', async () => {
  //     mockJwtService.verify.mockImplementation(() => {
  //       throw new Error('Invalid token');
  //     });
  //
  //     await expect(service.refreshToken('invalid-token')).rejects.toThrow(
  //       UnauthorizedException,
  //     );
  //   });
  // });
});

import { Test, TestingModule } from '@nestjs/testing';
import { AuthService } from './auth.service';
import { PrismaService } from '../prisma/prisma.service';
import { JwtService } from '@nestjs/jwt';
import { BadRequestException, UnauthorizedException } from '@nestjs/common';
import * as bcrypt from 'bcrypt';

describe('AuthService', () => {
  let service: AuthService;
  let prismaService: PrismaService;
  let jwtService: JwtService;

  const mockPrismaService = {
    user: {
      findUnique: jest.fn(),
      findFirst: jest.fn(),
      create: jest.fn(),
    },
    role: {
      findUnique: jest.fn(),
    },
  };

  const mockJwtService = {
    sign: jest.fn(),
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
    prismaService = module.get<PrismaService>(PrismaService);
    jwtService = module.get<JwtService>(JwtService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('findUserByUsername', () => {
    it('should find user by username', async () => {
      const mockUser = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.findUserByUsername('testuser');

      expect(result).toEqual(mockUser);
      expect(mockPrismaService.user.findUnique).toHaveBeenCalledWith({
        where: { username: 'testuser' },
      });
    });

    it('should return null if user not found', async () => {
      mockPrismaService.user.findUnique.mockResolvedValue(null);

      const result = await service.findUserByUsername('nonexistent');

      expect(result).toBeNull();
    });
  });

  describe('register', () => {
    it('should register a new user successfully', async () => {
      const registerDto = {
        username: 'newuser',
        email: 'new@example.com',
        password: 'SecurePass123!',
      };

      const mockRole = { id: 1, name: 'contestant' };
      const mockUser = { id: 1, username: 'newuser', email: 'new@example.com' };

      mockPrismaService.user.findFirst.mockResolvedValue(null);
      mockPrismaService.role.findUnique.mockResolvedValue(mockRole);
      mockPrismaService.user.create.mockResolvedValue(mockUser);

      const result = await service.register(registerDto);

      expect(result).toEqual({ user_id: 1 });
      expect(mockPrismaService.user.findFirst).toHaveBeenCalled();
      expect(mockPrismaService.user.create).toHaveBeenCalled();
    });

    it('should throw BadRequestException if username exists', async () => {
      const registerDto = {
        username: 'existing',
        email: 'new@example.com',
        password: 'SecurePass123!',
      };

      mockPrismaService.user.findFirst.mockResolvedValue({
        id: 1,
        username: 'existing',
      });

      await expect(service.register(registerDto)).rejects.toThrow(
        BadRequestException,
      );
    });

    it('should throw BadRequestException if email exists', async () => {
      const registerDto = {
        username: 'newuser',
        email: 'existing@example.com',
        password: 'SecurePass123!',
      };

      mockPrismaService.user.findFirst.mockResolvedValue({
        id: 1,
        email: 'existing@example.com',
      });

      await expect(service.register(registerDto)).rejects.toThrow(
        BadRequestException,
      );
    });
  });

  describe('login', () => {
    it('should return access token on successful login', async () => {
      const loginDto = {
        username: 'testuser',
        password: 'SecurePass123!',
      };

      const hashedPassword = await bcrypt.hash(loginDto.password, 12);
      const mockUser = {
        id: 1,
        username: 'testuser',
        password: hashedPassword,
        roleId: 1,
      };

      const mockRole = {
        id: 1,
        name: 'contestant',
        permissions: [{ name: 'submit_solution' }],
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);
      mockPrismaService.role.findUnique.mockResolvedValue(mockRole);
      mockJwtService.sign.mockReturnValue('mock-jwt-token');

      const result = await service.login(loginDto);

      expect(result).toHaveProperty('access_token');
      expect(result.access_token).toBe('mock-jwt-token');
      expect(mockJwtService.sign).toHaveBeenCalled();
    });

    it('should throw UnauthorizedException for invalid username', async () => {
      const loginDto = {
        username: 'nonexistent',
        password: 'password',
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

      const hashedPassword = await bcrypt.hash('CorrectPassword', 12);
      const mockUser = {
        id: 1,
        username: 'testuser',
        password: hashedPassword,
        roleId: 1,
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);

      await expect(service.login(loginDto)).rejects.toThrow(
        UnauthorizedException,
      );
    });
  });

  describe('verifyPermission', () => {
    it('should allow access with correct permission', async () => {
      const payload = { id: 1, username: 'testuser' };
      const mockUser = {
        id: 1,
        username: 'testuser',
        role: {
          name: 'admin',
          permissions: [{ name: 'manage_users' }],
        },
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.verifyPermission(payload, 'manage_users');

      expect(result.allowed).toBe(true);
      expect(result.user).toEqual(payload);
    });

    it('should deny access without permission', async () => {
      const payload = { id: 1, username: 'testuser' };
      const mockUser = {
        id: 1,
        username: 'testuser',
        role: {
          name: 'contestant',
          permissions: [{ name: 'submit_solution' }],
        },
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.verifyPermission(payload, 'manage_users');

      expect(result.allowed).toBe(false);
    });

    it('should allow wildcard permission', async () => {
      const payload = { id: 1, username: 'testuser' };
      const mockUser = {
        id: 1,
        username: 'testuser',
        role: {
          name: 'admin',
          permissions: [],
        },
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.verifyPermission(payload, '*');

      expect(result.allowed).toBe('*');
    });
  });

  describe('getProfile', () => {
    it('should return user profile without sensitive data', async () => {
      const mockUser = {
        id: 1,
        username: 'testuser',
        email: 'test@example.com',
        name: 'Test User',
        password: 'hashed-password',
        refreshToken: 'refresh-token',
        roleId: 1,
      };

      const mockRole = { id: 1, name: 'contestant' };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);
      mockPrismaService.role.findUnique.mockResolvedValue(mockRole);

      const result = await service.getProfile(1);

      expect(result).not.toHaveProperty('password');
      expect(result).not.toHaveProperty('refreshToken');
      expect(result.username).toBe('testuser');
      expect(result.role).toBe('contestant');
    });

    it('should throw UnauthorizedException if user not found', async () => {
      mockPrismaService.user.findUnique.mockResolvedValue(null);

      await expect(service.getProfile(999)).rejects.toThrow(
        UnauthorizedException,
      );
    });
  });

  describe('getPermissions', () => {
    it('should return user role and permissions', async () => {
      const mockUser = {
        id: 1,
        username: 'testuser',
        roleId: 1,
      };

      const mockRole = {
        id: 1,
        name: 'admin',
        permissions: [
          { name: 'manage_users' },
          { name: 'manage_problems' },
        ],
      };

      mockPrismaService.user.findUnique.mockResolvedValue(mockUser);
      mockPrismaService.role.findUnique.mockResolvedValue(mockRole);

      const result = await service.getPermissions(1);

      expect(result.role).toBe('admin');
      expect(result.permissions).toEqual(['manage_users', 'manage_problems']);
    });

    it('should throw UnauthorizedException if user not found', async () => {
      mockPrismaService.user.findUnique.mockResolvedValue(null);

      await expect(service.getPermissions(999)).rejects.toThrow(
        UnauthorizedException,
      );
    });
  });
});

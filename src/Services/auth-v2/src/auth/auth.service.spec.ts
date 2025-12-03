import { Test, TestingModule } from '@nestjs/testing';
import { AuthService } from './auth.service';
import { PrismaService } from '../prisma/prisma.service';
import { JwtService } from '@nestjs/jwt';

describe('AuthService', () => {
    let service: AuthService;
    let prisma: PrismaService;
    let jwt: JwtService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [
                AuthService,
                {
                    provide: PrismaService,
                    useValue: {
                        user: {
                            findUnique: jest.fn(),
                            findFirst: jest.fn(),
                            create: jest.fn(),
                        },
                        role: {
                            findUnique: jest.fn(),
                        },
                    },
                },
                {
                    provide: JwtService,
                    useValue: {
                        sign: jest.fn(),
                        verify: jest.fn(),
                    },
                },
            ],
        }).compile();

        service = module.get<AuthService>(AuthService);
        prisma = module.get<PrismaService>(PrismaService);
        jwt = module.get<JwtService>(JwtService);
    });

    it('should be defined', () => {
        expect(service).toBeDefined();
    });

    describe('verifyPermission', () => {
        it('should return allowed: true if user has permission', async () => {
            const mockUser = {
                id: 1,
                role: {
                    permissions: [{ name: 'test_perm' }],
                },
            };
            (prisma.user.findUnique as jest.Mock).mockResolvedValue(mockUser);

            const result = await service.verifyPermission({ id: 1 }, 'test_perm');
            expect(result.allowed).toBe(true);
        });

        it('should return allowed: false if user does not have permission', async () => {
            const mockUser = {
                id: 1,
                role: {
                    permissions: [{ name: 'other_perm' }],
                },
            };
            (prisma.user.findUnique as jest.Mock).mockResolvedValue(mockUser);

            const result = await service.verifyPermission({ id: 1 }, 'test_perm');
            expect(result.allowed).toBe(false);
        });
    });
});

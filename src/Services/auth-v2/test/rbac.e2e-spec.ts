import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication } from '@nestjs/common';
import request from 'supertest';
import { AppModule } from '../src/app.module';
import { PrismaService } from '../src/prisma/prisma.service';
import { seedDatabase } from '../prisma/seed';
import * as dotenv from 'dotenv';
import * as bcrypt from 'bcrypt';

dotenv.config();

describe('RBAC (e2e)', () => {
    let app: INestApplication;
    let prisma: PrismaService;
    let adminToken: string;
    let contestantToken: string;

    beforeAll(async () => {
        const moduleFixture: TestingModule = await Test.createTestingModule({
            imports: [AppModule],
        }).compile();

        app = moduleFixture.createNestApplication();
        await app.init();

        prisma = app.get<PrismaService>(PrismaService);

        // Ensure DB connection
        await prisma.$connect();

        // Seed database with roles and permissions from seed.ts
        await seedDatabase(prisma);

        // Cleanup test users
        await prisma.user.deleteMany({ where: { username: { in: ['admin_e2e', 'contestant_e2e'] } } });

        // Create Admin user
        await request(app.getHttpServer())
            .post('/auth/register')
            .send({ username: 'admin_e2e', password: 'password', email: 'admin_e2e@test.com' });

        const adminUser = await prisma.user.findUnique({ where: { username: 'admin_e2e' } });
        
        if (!adminUser) {
            throw new Error('Admin user not created');
        }

        // Get admin role
        const adminRole = await prisma.role.findUnique({ where: { name: 'admin' } });
        
        if (!adminRole) {
            throw new Error('Admin role not found');
        }

        // Update user to admin role BEFORE logging in
        await prisma.user.update({
            where: { id: adminUser.id },
            data: { roleId: adminRole.id }
        });

        // Now login with admin role already assigned
        const adminLogin = await request(app.getHttpServer())
            .post('/auth/login')
            .send({ username: 'admin_e2e', password: 'password' });
        adminToken = adminLogin.body.access_token;

        // Create Contestant (already has contestant role by default)
        await request(app.getHttpServer())
            .post('/auth/register')
            .send({ username: 'contestant_e2e', password: 'password', email: 'contestant_e2e@test.com' });

        const contestantLogin = await request(app.getHttpServer())
            .post('/auth/login')
            .send({ username: 'contestant_e2e', password: 'password' });
        contestantToken = contestantLogin.body.access_token;
    }, 60000); // Increase timeout to 60 seconds

    afterAll(async () => {
        await prisma.user.deleteMany({ where: { username: { in: ['admin_e2e', 'contestant_e2e'] } } });
        await prisma.$disconnect();
        await app.close();
    });

    it('/auth/verify (POST) - Admin should have manage_users permission', () => {
        return request(app.getHttpServer())
            .post('/auth/verify')
            .set('Authorization', `Bearer ${adminToken}`)
            .send({ permission: 'manage_users' })
            .expect(201)
            .expect((res) => {
                expect(res.body.allowed).toBe(true);
            });
    });

    it('/auth/verify (POST) - Contestant should NOT have manage_users permission', () => {
        return request(app.getHttpServer())
            .post('/auth/verify')
            .set('Authorization', `Bearer ${contestantToken}`)
            .send({ permission: 'manage_users' })
            .expect(201)
            .expect((res) => {
                expect(res.body.allowed).toBe(false);
            });
    });

    it('/auth/verify (POST) - Contestant should have submit_code permission', () => {
        return request(app.getHttpServer())
            .post('/auth/verify')
            .set('Authorization', `Bearer ${contestantToken}`)
            .send({ permission: 'submit_code' })
            .expect(201)
            .expect((res) => {
                expect(res.body.allowed).toBe(true);
            });
    });

    it('/auth/verify (POST) - Invalid token should return 401', () => {
        return request(app.getHttpServer())
            .post('/auth/verify')
            .set('Authorization', 'Bearer invalid_token')
            .send({ permission: 'submit_code' })
            .expect(401);
    });
});

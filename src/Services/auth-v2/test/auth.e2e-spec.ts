import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, ValidationPipe } from '@nestjs/common';
import request from 'supertest';
import { AppModule } from './../src/app.module';
import { PrismaService } from './../src/prisma/prisma.service';
import { seedDatabase } from './../prisma/seed';

describe('AuthController (e2e)', () => {
  let app: INestApplication;
  let prisma: PrismaService;

  beforeAll(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    app.useGlobalPipes(new ValidationPipe());
    prisma = moduleFixture.get<PrismaService>(PrismaService);
    await app.init();

    // Seed database with roles and permissions from seed.ts
    await seedDatabase(prisma);
  });

  afterAll(async () => {
    await prisma.$disconnect();
    await app.close();
  });

  beforeEach(async () => {
    // Clean only users, keep roles and permissions from seed
    await prisma.user.deleteMany({
      where: {
        username: { not: 'admin' } // Keep the admin user from seed
      }
    });
  });

  describe('/auth/register (POST)', () => {
    it('should register a new user', async () => {
      const response = await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'testuser',
          password: 'SecurePass123!',
        })
        .expect(201);

      expect(response.body).toHaveProperty('user_id');
      expect(typeof response.body.user_id).toBe('number');
    });

    it('should fail with duplicate username', async () => {
      // Register first user
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'user1@example.com',
          username: 'testuser',
          password: 'Pass123!',
        });

      // Try to register with same username
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'user2@example.com',
          username: 'testuser',
          password: 'Pass123!',
        })
        .expect(400);
    });

    it('should fail with duplicate email', async () => {
      // Register first user
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'user1',
          password: 'Pass123!',
        });

      // Try to register with same email
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'user2',
          password: 'Pass123!',
        })
        .expect(400);
    });
  });

  describe('/auth/login (POST)', () => {
    beforeEach(async () => {
      // Create a test user
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'testuser',
          password: 'SecurePass123!',
        });
    });

    it('should login successfully with username', async () => {
      const response = await request(app.getHttpServer())
        .post('/auth/login')
        .send({
          username: 'testuser',
          password: 'SecurePass123!',
        })
        .expect(201);

      expect(response.body).toHaveProperty('access_token');
      expect(typeof response.body.access_token).toBe('string');
    });

    it('should fail with wrong password', () => {
      return request(app.getHttpServer())
        .post('/auth/login')
        .send({
          username: 'testuser',
          password: 'WrongPassword',
        })
        .expect(401);
    });

    it('should fail with non-existent user', () => {
      return request(app.getHttpServer())
        .post('/auth/login')
        .send({
          username: 'nonexistent',
          password: 'SecurePass123!',
        })
        .expect(401);
    });
  });

  describe('/auth/profile (GET)', () => {
    let accessToken: string;

    beforeEach(async () => {
      // Register and login
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'testuser',
          password: 'SecurePass123!',
        });

      const loginRes = await request(app.getHttpServer())
        .post('/auth/login')
        .send({
          username: 'testuser',
          password: 'SecurePass123!',
        });

      accessToken = loginRes.body.access_token;
    });

    it('should get current user profile with valid token', async () => {
      const response = await request(app.getHttpServer())
        .get('/auth/profile')
        .set('Authorization', `Bearer ${accessToken}`)
        .expect(200);

      expect(response.body.email).toBe('test@example.com');
      expect(response.body.username).toBe('testuser');
      expect(response.body).not.toHaveProperty('password');
    });

    it('should fail without token', () => {
      return request(app.getHttpServer())
        .get('/auth/profile')
        .expect(401);
    });

    it('should fail with invalid token', () => {
      return request(app.getHttpServer())
        .get('/auth/profile')
        .set('Authorization', 'Bearer invalid-token')
        .expect(401);
    });
  });

  describe('/auth/validate (POST)', () => {
    let accessToken: string;

    beforeEach(async () => {
      await request(app.getHttpServer())
        .post('/auth/register')
        .send({
          email: 'test@example.com',
          username: 'testuser',
          password: 'SecurePass123!',
        });

      const loginRes = await request(app.getHttpServer())
        .post('/auth/login')
        .send({
          username: 'testuser',
          password: 'SecurePass123!',
        });

      accessToken = loginRes.body.access_token;
    });

    it('should validate token successfully', async () => {
      const response = await request(app.getHttpServer())
        .post('/auth/validate')
        .set('Authorization', `Bearer ${accessToken}`)
        .expect(201);

      expect(response.body).toHaveProperty('id');
      expect(response.body).toHaveProperty('username');
      expect(response.body.username).toBe('testuser');
    });

    it('should fail with invalid token', () => {
      return request(app.getHttpServer())
        .post('/auth/validate')
        .set('Authorization', 'Bearer invalid-token')
        .expect(401);
    });
  });

  describe('/auth/health (GET)', () => {
    it('should return health status', async () => {
      const response = await request(app.getHttpServer())
        .get('/auth/health')
        .expect(200);

      expect(response.body).toHaveProperty('status');
      expect(response.body.status).toBe('ok');
      expect(response.body).toHaveProperty('date');
    });
  });
});

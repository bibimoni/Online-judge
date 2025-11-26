import { Injectable, UnauthorizedException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { UsersService } from '../users/users.service';
import * as bcrypt from 'bcrypt';
import { PrismaService } from 'src/prisma/prisma.service';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtService: JwtService,
    private prisma: PrismaService
  ) { }
  async findUserByUsername(username: string) {
    return this.prisma.user.findUnique({
      where: { username }
    });
  }

  async validateUser(username: string, password: string) {
    const user = await this.findUserByUsername(username);
    if (user && (await bcrypt.compare(password, user.password))) {
      const { password, ...result } = user;
      return result;
    }
    return null;
  }

  async login(username: string, password: string) {
    const user = await this.validateUser(username, password);
    if (!user) {
      throw new UnauthorizedException('Invalid credentials');
    }
    const payload = { username: user.username, sub: user.id };
    return this.jwtService.sign(payload)
  }

  async validateToken(token: string) {
    try {
      const decoded = this.jwtService.verify(token);

      console.log('Decoded token:', decoded);
      return {
        code: 0,
        user: {
          username: decoded.username,
          id: decoded.sub,
        }
      }
    } catch (error) {
      return {
        code: 1,
        message: 'Invalid token'
      };
    }
  }

  async refreshToken(refreshToken: string) {
    try {
      const decoded = this.jwtService.verify(refreshToken);
      const user = await this.usersService.findByEmail(decoded.email);

      if (!user) {
        throw new UnauthorizedException('User not found');
      }

      const payload = { username: user.username, sub: user.id };
      return this.jwtService.sign(payload);
    } catch (error) {
      throw new UnauthorizedException('Invalid refresh token');
    }
  }
}

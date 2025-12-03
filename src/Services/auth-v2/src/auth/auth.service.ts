import { Injectable, UnauthorizedException, BadRequestException } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
import { PrismaService } from 'src/prisma/prisma.service';
import { RegisterDto } from './dto/register.dto';
import { LoginDto } from './dto/login.dto';
import { config } from 'src/config/config';

@Injectable()
export class AuthService {
  constructor(
    public jwtService: JwtService,
    private prisma: PrismaService
  ) { }

  async findUserByUsername(username: string) {
    return this.prisma.user.findUnique({ where: { username } });
  }

  async register(registerDto: RegisterDto) {
    const { username, email, password } = registerDto;
    const existingUser = await this.prisma.user.findFirst({
      where: { OR: [{ username }, { email }] }
    });
    if (existingUser) {
      throw new BadRequestException('Username or email already exists');
    }
    const hashedPassword = await bcrypt.hash(password, 12);
    const defaultRole = await this.prisma.role.findUnique({ where: { name: 'contestant' } });
    if (!defaultRole) throw new BadRequestException('Default role not found');
    const user = await this.prisma.user.create({
      data: {
        username,
        email,
        password: hashedPassword,
        roleId: defaultRole.id,
        name: username,
      },
    });
    return { user_id: user.id };
  }

  async login(loginDto: LoginDto) {
    const user = await this.findUserByUsername(loginDto.username);
    if (!user || !(await bcrypt.compare(loginDto.password, user.password))) {
      throw new UnauthorizedException('Invalid credentials');
    }
    const role = await this.prisma.role.findUnique({
      where: { id: user.roleId },
      include: { permissions: true }
    });
    const permissions = role ? role.permissions.map(p => p.name) : [];
    const payload = {
      id: user.id,
      username: user.username,
      role: role ? role.name : null,
      permissions,
    };
    const accessToken = this.jwtService.sign(payload, { expiresIn: config.jwtExpiresIn });
    return { access_token: accessToken };
  }

  async verifyPermission(payload: any, requiredPermission: string) {
    try {
      const user = await this.prisma.user.findUnique({
        where: { id: payload.id },
        include: { role: { include: { permissions: true } } }
      });

      if (!user || !user.role) return { allowed: false };

      if (requiredPermission == "*") {
        return { allowed: "*", user: payload }
      }
      const hasPermission = user.role.permissions.some(p => p.name === requiredPermission);
      return { allowed: hasPermission, user: payload };
    } catch (e) {
      throw new UnauthorizedException('Invalid token');
    }
  }

  async getProfile(userId: number) {
    const user = await this.prisma.user.findUnique({ where: { id: userId } });
    if (!user) throw new UnauthorizedException('User not found');
    const role = user.roleId ? await this.prisma.role.findUnique({ where: { id: user.roleId } }) : null;
    const { password, refreshToken, ...result } = user;
    return { ...result, role: role ? role.name : null };
  }

  async getPermissions(userId: number) {
    const user = await this.prisma.user.findUnique({ where: { id: userId } });
    if (!user) throw new UnauthorizedException('User not found');
    const role = user.roleId ? await this.prisma.role.findUnique({
      where: { id: user.roleId },
      include: { permissions: true }
    }) : null;
    const permissions = role ? role.permissions.map(p => p.name) : [];
    return { role: role ? role.name : null, permissions };
  }
}

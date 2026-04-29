import { Injectable } from '@nestjs/common';
import { PrismaService } from '../prisma/prisma.service';
import { CreateUserDto } from './dto/create-user.dto';
import * as bcrypt from 'bcrypt';

@Injectable()
export class UsersService {
  constructor(private prisma: PrismaService) { }

  async create(createUserDto: CreateUserDto) {
    const { username, password, name } = createUserDto;
    const hashedPassword = await bcrypt.hash(password, 10);
    return this.prisma.user.create({
      data: {
        username,
        password: hashedPassword,
        name,
      },
    });
  }

  async findByEmail(email: string) {
    return this.prisma.user.findUnique({ where: { email } });
  }

  async findById(user_id: number) { // Changed from string to number
    return this.prisma.user.findUnique({ where: { id: user_id } });
  }


  async findRoleById(roleId: number) {
    const role = await this.prisma.role.findUnique({
      where: { id: roleId },
      include: { permissions: true }
    });
    return role;
  }

}

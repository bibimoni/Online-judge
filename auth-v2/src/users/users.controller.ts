import { Controller, Get, UseGuards, Request } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBearerAuth } from '@nestjs/swagger';
import { JwtAuthGuard } from '../auth/jwt.strategy';
import { UsersService } from './users.service';

@ApiTags('Users')
@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) { }

  // Remove direct user creation, use /auth/register instead

  // Remove find by email for privacy and security

  @Get('profile')
  @UseGuards(JwtAuthGuard)
  @ApiOperation({ summary: 'Get current authenticated user profile' })
  @ApiBearerAuth('JWT-auth')
  @ApiResponse({
    status: 200,
    description: 'Current user info',
    schema: {
      type: 'object',
      properties: {
        id: { type: 'number' },
        username: { type: 'string' },
        email: { type: 'string' },
        name: { type: 'string' },
        role: { type: 'string' }
      }
    }
  })
  async getProfile(@Request() req) {
    const dbUser = await this.usersService.findById(req.user.id);
    const role = await this.usersService.findRoleById(dbUser.roleId);

    return {
      id: dbUser.id,
      username: dbUser.username,
      email: dbUser.email,
      name: dbUser.name,
      role: role ? role.name : null
    };
  }
}

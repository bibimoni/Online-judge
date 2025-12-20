import { Controller, Post, Body, Get, UseGuards, Request, Headers, UnauthorizedException } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse, ApiBody, ApiBearerAuth } from '@nestjs/swagger';
import { AuthService } from './auth.service';
import { JwtAuthGuard } from './jwt.strategy';
import { LoginDto } from './dto/login.dto';
import { RegisterDto } from './dto/register.dto';

import { PermissionsGuard } from './permissions.guard';
import { Permissions } from './permissions.decorator';
import { PermissionVerifyDto } from './dto/permission-verify.dto';

@ApiTags('Authentication')
@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) { }

  @Post('register')
  @ApiOperation({ summary: 'User registration' })
  @ApiBody({ type: RegisterDto })
  @ApiResponse({
    status: 201,
    description: 'User registered',
    schema: {
      type: 'object',
      properties: {
        user_id: { type: 'number', example: 1 }
      }
    }
  })
  async register(@Body() registerDto: RegisterDto) {
    return await this.authService.register(registerDto);
  }

  @Post('login')
  @ApiOperation({ summary: 'User login' })
  @ApiBody({ type: LoginDto })
  @ApiResponse({
    status: 200,
    description: 'Login successful',
    schema: {
      type: 'object',
      properties: {
        access_token: { type: 'string' }
      }
    }
  })
  async login(@Body() loginDto: LoginDto) {
    return await this.authService.login(loginDto);
  }

  @UseGuards(JwtAuthGuard)
  @Post('verify')
  @ApiOperation({ summary: 'Check if user have the required permission' })
  @ApiBody({
    type: PermissionVerifyDto
  })
  @ApiBearerAuth("JWT-auth")
  async verify(@Request() req: any, @Body() body: PermissionVerifyDto) {
    if (!req.user) {
      throw new UnauthorizedException()
    }
    
    return await this.authService.verifyPermission(req.user, body.permission);
  }

  @UseGuards(JwtAuthGuard)
  @Get('profile')
  @ApiOperation({ summary: 'Get current user profile' })
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
        role: { type: 'string' },
        permissions: {
          type: 'array',
          items: { type: 'string' }
        }
      }
    }
  })
  async getProfile(@Request() req) {
    return await this.authService.getProfile(req.user.id);
  }

  @UseGuards(JwtAuthGuard, PermissionsGuard)
  @Permissions('view_profile')
  @Get('permissions')
  @ApiOperation({ summary: 'Get current user permissions and role' })
  @ApiBearerAuth('JWT-auth')
  @ApiResponse({
    status: 200,
    description: 'User permissions and role',
    schema: {
      type: 'object',
      properties: {
        role: { type: 'string' },
        permissions: {
          type: 'array',
          items: { type: 'string' }
        }
      }
    }
  })
  async getPermissions(@Request() req: any) {
    return await this.authService.getPermissions(req.user.id);
  }

  @Get('health')
  @ApiOperation({ summary: 'Health check' })
  @ApiResponse({ status: 200, description: 'Service is healthy' })
  getHealth() {
    return { status: 'ok', date: new Date().toISOString() };
  }

  @UseGuards(JwtAuthGuard)
  @Post('validate')
  @ApiOperation({ summary: 'Validate JWT token' })
  @ApiBearerAuth('JWT-auth')
  @ApiResponse({
    status: 200,
    description: 'Token is valid',
    schema: {
      type: 'object',
      properties: {
        success: { type: 'boolean', example: true },
        id: { type: 'number' },
        username: { type: 'string' },
        role: { type: 'string' },
        permissions: { type: 'array', items: { type: 'string' } }
      }
    }
  })
  @ApiResponse({ status: 401, description: 'Invalid token' })
  async validate(@Request() req: any) {
    return req.user;
  }
}

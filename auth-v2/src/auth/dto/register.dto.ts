import { IsString, IsEmail, MinLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class RegisterDto {
  @ApiProperty({ description: 'Username', example: 'john_doe' })
  @IsString()
  username: string;

  @ApiProperty({ description: 'Email', example: 'john@example.com' })
  @IsEmail()
  email: string;

  @ApiProperty({ description: 'Password', example: 'password123', minLength: 6 })
  @IsString()
  @MinLength(6)
  password: string;
}

import { IsEmail, IsString, MinLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class LoginDto {
  @ApiProperty({ 
    description: 'Username for the user',
    example: 'john_doe'
  })
  @IsString()
  username: string;

  @ApiProperty({ 
    description: 'User password (minimum 6 characters)',
    example: 'password123',
    minLength: 6
  })
  @IsString()
  @MinLength(6)
  password: string;
}

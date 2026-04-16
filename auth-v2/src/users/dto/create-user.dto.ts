import { IsEmail, IsString, MinLength } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class CreateUserDto {
  @ApiProperty({ 
    description: 'Username for the user',
    example: 'john_doe'
  })
  username: string;

  @ApiProperty({ 
    description: 'User password (minimum 6 characters)',
    example: 'password123',
    minLength: 6
  })
  @IsString()
  @MinLength(6)
  password: string;

  @ApiProperty({ 
    description: 'User full name',
    example: 'John Doe'
  })
  @IsString()
  name: string;
}
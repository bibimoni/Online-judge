import { IsString } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';

export class PermissionVerifyDto {
  @ApiProperty({
    description: 'Permission to check',
    example: 'submit_code'
  })
  @IsString()
  permission: string;
}

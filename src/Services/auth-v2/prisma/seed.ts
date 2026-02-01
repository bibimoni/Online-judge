import { PrismaClient } from '../generated/client';
import * as bcrypt from 'bcrypt';

// Export the seed function so it can be reused in tests
export async function seedDatabase(prisma: PrismaClient) {
  console.log('Starting seed...');
  const permissions = [
    { name: 'login', description: 'Can login to the system' },
    { name: 'register', description: 'Can register a new account' },

    { name: 'view_problem', description: 'Can view problems' },
    { name: 'create_problem', description: 'Can create new problems' },
    { name: 'edit_problem', description: 'Can edit existing problems' },
    { name: 'delete_problem', description: 'Can delete problems' },
    { name: 'view_hidden_problem', description: 'Can view hidden problems' },

    { name: 'submit_code', description: 'Can submit code for problems' },
    { name: 'view_submission', description: 'Can view submissions' },
    { name: 'rejudge_submission', description: 'Can rejudge submissions' },

    { name: 'view_contest', description: 'Can view contests' },
    { name: 'register_contest', description: 'Can register for contests' },
    { name: 'create_contest', description: 'Can create new contests' },
    { name: 'edit_contest', description: 'Can edit existing contests' },
    { name: 'delete_contest', description: 'Can delete contests' },
    { name: 'manage_contest', description: 'Can manage contests (add problems, etc.)' },

    { name: 'view_profile', description: 'Can view user profiles' },
    { name: 'edit_profile', description: 'Can edit own profile' },
    { name: 'manage_users', description: 'Can manage other users' },
    { name: 'admin', description: 'Super admin permission' },
  ];

  for (const p of permissions) {
    await prisma.permission.upsert({
      where: { name: p.name },
      update: { description: p.description },
      create: { name: p.name, description: p.description },
    });
  }

  const getPerms = async (names: string[]) => {
    const perms = await prisma.permission.findMany({ where: { name: { in: names } } });
    return perms.map(p => ({ id: p.id }));
  };

  const allPermissions = await prisma.permission.findMany();

  const roles = [
    {
      name: 'contestant',
      description: 'Regular user who can participate in contests',
      permissions: [
        'login', 'register',
        'view_problem', 'submit_code', 'view_submission',
        'view_contest', 'register_contest',
        'view_profile', 'edit_profile'
      ]
    },
    {
      name: 'problem_setter',
      description: 'User who can create problems and contests',
      permissions: [
        'login', 'register',
        'view_problem', 'create_problem', 'edit_problem', 'view_hidden_problem',
        'submit_code', 'view_submission', 'rejudge_submission',
        'view_contest', 'register_contest', 'create_contest', 'edit_contest', 'manage_contest',
        'view_profile', 'edit_profile'
      ]
    },
    {
      name: 'admin',
      description: 'Administrator with full access',
      permissions: allPermissions.map(p => p.name)
    }
  ];

  for (const r of roles) {
    const permIds = await getPerms(r.permissions);
    
    await prisma.role.upsert({
      where: { name: r.name },
      update: {
        description: r.description,
        permissions: { set: permIds }
      },
      create: {
        name: r.name,
        description: r.description,
        permissions: { connect: permIds }
      }
    });
  }

  // Create Admin User
  const password = await bcrypt.hash('bkacbkac', 10);
  await prisma.user.upsert({
    where: { username: 'admin' },
    update: {
      password: password,
      role: { connect: { name: 'admin' } },
    },
    create: {
      username: 'admin',
      name: 'Admin User',
      email: 'admin@example.com',
      password: password,
      role: { connect: { name: 'admin' } },
    },
  });
  console.log('Admin user seeded');
}

// Main function for running seed from CLI
async function main() {
  const prisma = new PrismaClient();
  try {
    await seedDatabase(prisma);
  } finally {
    await prisma.$disconnect();
  }
}

// Only run main if this file is executed directly
if (require.main === module) {
  main();
}

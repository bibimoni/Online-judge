import { PrismaClient } from '../generated/client'

const prisma = new PrismaClient();
async function main() {
  const permCreateContest = await prisma.permission.upsert({
    where: { name: 'contest:create' },
    update: {},
    create: { name: 'contest:create' },
  });
  const adminRole = await prisma.role.upsert({
    where: { name: 'ADMIN' },
    update: {},
    create: {
      name: 'ADMIN',
      permissions: { connect: [{ name: 'contest:create' }] },
    },
  });
  const userRole = await prisma.role.upsert({
    where: { name: 'USER' },
    update: {},
    create: { name: 'USER' },
  });
  console.log({ adminRole, userRole });
}
main()
  .catch((e) => console.error(e))
  .finally(async () => await prisma.$disconnect());

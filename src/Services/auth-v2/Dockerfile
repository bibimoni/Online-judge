FROM node:20
WORKDIR /app
COPY package*.json ./
COPY prisma ./prisma
RUN npm ci
RUN npx prisma generate
COPY . .
EXPOSE 50051
CMD ["npm", "run", "start:dev"]

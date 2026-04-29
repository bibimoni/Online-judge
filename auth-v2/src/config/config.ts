export class Config {
  private static instance: Config;
  private _config: any = {};

  private constructor() { this.load() }

  static getInstance(): Config {
    if (!Config.instance) {
      Config.instance = new Config();
    }
    return Config.instance;
  }

  load() {
    this._config = {
      port: parseInt(process.env.PORT!) || 50051,
      jwtSecret: process.env.AUTH_JWT_SECRET || 'secret',
      jwtExpiresIn: process.env.JWT_EXPIRES_IN || '7d',
      isDevelopment: process.env.NODE_ENV || 'development'
    };
    this._config.databaseUri = process.env.AUTH_DATABASE_URL || `postgresql://postgres:postgres@localhost:5432/auth_db`

    return this;
  }

  get port() { return this._config.port }
  get jwtSecret() { return this._config.jwtSecret }
  get jwtExpiresIn() { return this._config.jwtExpiresIn }
  get isDevelopment() { return this._config.isDevelopment }
  get databaseUri() { return this._config.databaseUri }

  get(key: string) {
    return this._config[key];
  }

  getAll() {
    return { ...this._config };
  }
}

export const config = Config.getInstance();


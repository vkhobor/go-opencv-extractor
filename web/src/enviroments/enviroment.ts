import { Enviroment } from './enviroment.d';

export default {
    production: false,
    api: import.meta.env['NG_APP_API_URL'] ?? '/api',
} satisfies Enviroment;

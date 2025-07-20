import env from '../../../environments/environment';
import { Api } from '../../../api/Api';
export const baseUrl = env.api.endsWith('/api')
    ? env.api.slice(0, -4)
    : env.api;

export const client = new Api({ baseUrl });

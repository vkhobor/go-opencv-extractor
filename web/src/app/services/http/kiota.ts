import { FetchRequestAdapter } from '@microsoft/kiota-http-fetchlibrary';
import { createApiClient } from '../../../api/apiClient';
import { AnonymousAuthenticationProvider } from '@microsoft/kiota-abstractions';
import env from '../../../enviroments/enviroment';

const authProvider = new AnonymousAuthenticationProvider();
const adapter = new FetchRequestAdapter(authProvider);
env.api.endsWith('/api')
    ? (adapter.baseUrl = env.api.slice(0, -4))
    : (adapter.baseUrl = env.api);
export const client = createApiClient(adapter);

import 'normalize.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/700.css';
import { GlobalStyle } from './styles/Global';

import type {} from 'styled-components/cssprop';
import { Dashboard } from './components/Dashboard';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/Theme';
import { QueryClientProvider, QueryClient } from 'react-query';
const queryClient = new QueryClient();
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <GlobalStyle />
        <Dashboard />
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;

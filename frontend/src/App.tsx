import 'normalize.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/700.css';
import { GlobalStyle } from './styles/Global';

import type {} from 'styled-components/cssprop';
import { Dashboard } from './components/Dashboard';
import { ThemeProvider } from 'styled-components';
import { theme } from './styles/Theme';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <GlobalStyle />
      <Dashboard />
    </ThemeProvider>
  );
}

export default App;

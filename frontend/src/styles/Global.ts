import { createGlobalStyle } from 'styled-components';

const GlobalStyle = createGlobalStyle<{}>`
    html {
    box-sizing: border-box;
    font-size: 62.5%;
    }
    *, *:before, *:after {
    box-sizing: inherit;
    }
    body{
        font-family: 'Roboto', sans-serif;
        font-size: 1.6rem;
    }
`;

export { GlobalStyle };

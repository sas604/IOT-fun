const defautlTheme = {
  light: '#F5F6F4',
  neutral: '#b4b8abff',
  secondary: '#153243',
  tertiary: '#284b63ff',
  tertiaryLight: '#D4E3ED',
  highlight: '#fa8334ff',
  primary: '#e55934ff',
};
const sizing = {
  xs: '0.5rem',
  sm: '1rem',
  md: '1.5rem',
  lg: '2rem',
  xl: '2.5rem',
  '2xl': '3rem',
};

const border = {
  colored: (color: string) => `1px solid ${color}`,
  r: '4px',
  'r-lg': '8px',
};

const theme = {
  color: defautlTheme,
  ...sizing,
  border,
};

export { theme };

import { styled } from 'styled-components';

import { SwitchList } from './SwitchList';
import { useQuery } from 'react-query';

function Dashboard() {
  const { data, isLoading, isError } = useQuery('switches', async () => {
    const res = await fetch('/api/switches');
    if (!res.ok) {
      throw new Error('Network response was not ok');
    }
    return res.json();
  });

  if (isLoading) {
    return (
      <DashBoardStyle>
        <h1>Dashboard</h1>
        <p css="color:white">Loading....</p>
      </DashBoardStyle>
    );
  }
  return (
    <DashBoardStyle>
      <h1>Dashboard</h1>
      <div></div>
      <SwitchList switches={data} />
    </DashBoardStyle>
  );
}
const DashBoardStyle = styled.div`
  padding: 0 20px 20px;
  background-color: ${({ theme }) => theme.color.tertiary};
  min-height: 100vh;
  h1 {
    margin-top: 0;
    padding-top: ${({ theme }) => theme.lg};
    color: ${({ theme }) => theme.color.light};
  }
`;

export { Dashboard };

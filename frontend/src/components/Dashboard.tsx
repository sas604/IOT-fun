import { styled } from 'styled-components';
import { SwitchList } from './SwitchList';
import { SwitchProps } from './Switch';
import { useQuery } from '@tanstack/react-query';

async function fetchSwitches(): Promise<SwitchProps[]> {
  const res = await fetch('/api/switches');
  if (!res.ok) {
    throw new Error('Network response was not ok');
  }
  return res.json();
}

function Dashboard() {
  const { isError, error, isLoading, data } = useQuery({
    queryKey: ['switches'],
    queryFn: fetchSwitches,
  });

  if (isError) {
    return (
      <DashBoardStyle>
        <h1>Dashboard</h1>
        <p css="color:white">{error.message}</p>
      </DashBoardStyle>
    );
  }

  if (isLoading) {
    return (
      <DashBoardStyle>
        <h1>Dashboard</h1>
        <p css="color:white">Loading....</p>
      </DashBoardStyle>
    );
  }
  if (data) {
    return (
      <DashBoardStyle>
        <h1>Dashboard</h1>
        <div></div>
        <SwitchList switches={data} />
      </DashBoardStyle>
    );
  }
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

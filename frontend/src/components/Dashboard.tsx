import { styled } from 'styled-components';
import { SwitchProps } from './Switch';
import { SwitchList } from './SwitchList';

const demoSwitch: SwitchProps[] = [
  {
    state: true,
    autoControl: true,
    value: 36,
    target: 25,
    measurment: 'Temperature',
    unit: 'C',
  },
  {
    state: false,
    autoControl: true,
    value: 75,
    target: 90,
    measurment: 'Humidity',
    unit: '%',
  },
];

function Dashboard() {
  return (
    <DashBoardStyle>
      <h1>Dashboard</h1>
      <div></div>
      <SwitchList switches={demoSwitch} />
    </DashBoardStyle>
  );
}
const DashBoardStyle = styled.div`
  padding: 0 20px;
  background-color: ${({ theme }) => theme.color.tertiary};
  min-height: 100vh;
  h1 {
    margin-top: 0;
    padding-top: ${({ theme }) => theme.lg};
    color: ${({ theme }) => theme.color.light};
  }
`;

export { Dashboard };

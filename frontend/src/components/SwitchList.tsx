import { styled } from 'styled-components';
import { Switch, SwitchProps } from './Switch';

export type SwitchList = {
  switches: SwitchProps[];
};

const SwitchList: React.FC<SwitchList> = ({ switches }) => {
  return (
    <SwitchListStyles>
      {switches.map(
        ({ state, measurment, value, target, autoControl, unit, id }, idx) => (
          <li key={measurment}>
            <Switch
              state={state}
              measurment={measurment}
              target={target}
              value={value}
              autoControl={autoControl}
              unit={unit}
              id={id}
              idx={idx}
            ></Switch>
          </li>
        )
      )}
    </SwitchListStyles>
  );
};

const SwitchListStyles = styled.ul`
  list-style: none;
  padding: 0;
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.lg};
  > * {
    flex: 1 1 300px;
  }
`;

export { SwitchList };

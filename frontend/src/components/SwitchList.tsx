import { styled } from 'styled-components';
import { Switch, SwitchProps } from './Switch';

export type SwitchList = {
  switches: SwitchProps[];
};

const SwitchList: React.FC<SwitchList> = ({ switches }) => {
  return (
    <SwitchListStyles>
      {switches.map((switchData) => (
          <li key={switchData.name}>
            <Switch
              switchData = {switchData}
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

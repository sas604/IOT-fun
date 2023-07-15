import { styled } from 'styled-components';
import { ImSwitch } from 'react-icons/im';

export type SwitchProps = {
  state: 'off' | 'on';
  measurment: string;
  value: number;
  autoControl: boolean;
  target?: number;
  unit: string;
};

const Switch: React.FC<SwitchProps> = ({
  state,
  measurment,
  value,
  autoControl,
  target,
  unit,
}) => {
  return (
    <SwitchStyle>
      <h2>{measurment}</h2>
      <div>
        <div>
          <p>Current Value: </p>
          <p className="switch-value">
            {Math.round(value * 10) / 10}
            {unit}
          </p>
        </div>
        <SwitchToggleStyle $active={state === 'on'}>
          <button>
            <ImSwitch></ImSwitch>
          </button>
          <input
            id={measurment + '-switch'}
            onChange={(e) => console.log(e)}
            type="checkbox"
            checked={state === 'on'}
          />
        </SwitchToggleStyle>
      </div>

      <SwitchFooterStye>
        <div>
          <span css="font-weight: 700">Automation: </span>
          <span>{autoControl ? 'Enabled' : 'Disabled'}</span>
        </div>
        <div>
          <span css="font-weight: 700">Automation Target Value: </span>
          <span>
            {target}
            {unit}
          </span>
        </div>
      </SwitchFooterStye>
    </SwitchStyle>
  );
};

const SwitchStyle = styled.div`
  padding: ${({ theme }) => theme.lg};
  border-radius: ${({ theme }) => theme.border.r};
  background-color: ${({ theme }) => theme.color.light};
  border: ${({ theme }) => theme.border.colored(theme.color.secondary)};
  color: ${({ theme }) => theme.color.secondary};
  h2 {
    margin-top: 0;
    margin-bottom: ${({ theme }) => theme.xs};
  }

  h2 + div {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  p {
    margin: 0;
  }
  .switch-value {
    font-size: 2em;
    font-weight: 700;
  }
`;

const SwitchToggleStyle = styled.div<{ $active?: boolean }>`
  button {
    border-radius: 50%;
    font-size: 1em;
    color: ${({ theme, $active }) =>
      $active ? 'white' : theme.color.secondary};
    appearance: none;
    cursor: pointer;
    background-color: ${({ theme, $active }) =>
      $active ? theme.color.primary : theme.color.tertiaryLight};

    border: ${({ theme, $active }) =>
      $active ? 'none' : theme.border.colored(theme.color.secondary)};
    padding: 10px 8px;
    width: 2.5em;
    height: 2.5em;
  }
  input {
    position: absolute;
    width: 1px;
    height: 1px;
  }
`;

const SwitchFooterStye = styled.div`
  margin-top: ${({ theme }) => theme.sm};
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.sm};
  justify-content: space-between;
`;
export { Switch };

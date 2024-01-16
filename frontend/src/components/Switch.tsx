import { styled } from 'styled-components';
import { ImSwitch } from 'react-icons/im';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { setSwitch } from '../lib/setSwitch';

export interface SwitchProps  {
  state: 'off' | 'on';
  measurement: string;
  value: number;
  automation: boolean;
  schedule: boolean;
  interval: number;
  duration: number;
  maxValue: number;
  minValue: number;
  unit: string;
  name: string;

  
}


// const queryClient = useQueryClient();
// const mutation = useMutation(setSwitch, {
//   onMutate: async (newState) => {
//     // Cancel any outgoing refetches (so they don't overwrite our optimistic update)
//     console.log(newState);
//     await queryClient.cancelQueries('switches');

//     const previousState = queryClient.getQueryData<SwitchProps[]>('switches');
//     if (previousState) {
//       previousState[newState.idx].state = newState.state;

//       queryClient.setQueryData<SwitchProps[]>('switches', previousState);
//     }

//     // Return a context object with the snapshotted value
//     return { previousState };
//   },
//   onError: (err, newState, context) => {
//     if (context?.previousState) {
//       queryClient.setQueryData('switches', context.previousState);
//     }
//   },
//   onSettled: () => {
//     queryClient.invalidateQueries('switches');
//   },
// });
const Switch = ({switchData }: {switchData : SwitchProps}) => {
const {measurement, value, unit, state, automation, schedule, maxValue, minValue, duration, interval, name} = switchData
const queryClient = useQueryClient();
const mutation = useMutation({mutationFn: setSwitch, onMutate:async (newState) => {
  console.log(newState)
  await queryClient.cancelQueries({queryKey:['switches']})
  const prev = queryClient.getQueryData<SwitchProps[]>(['switches']);

  queryClient.setQueryData(['switches'], (old)=> console.log(old))

}});
  return (
    <SwitchStyle>
      <h2>{name} {measurement && " - " + measurement} </h2>
      <div>
         {value ? (<div>
          <p>Current Value: </p>
          <p className="switch-value">
            {Math.round(value * 10) / 10}
            {unit}
          </p>
        </div>) : (<div>
          <p>State: </p>
          <p className="switch-value">
            {state}
          </p>
        </div>)}
        <SwitchToggleStyle $active={state === 'on'}>
          <button onClick={(e) => mutation.mutate({id:"1",state:"off"})}>
            <ImSwitch></ImSwitch>
          </button>
          <input
            id={measurement + '-switch'}
            onChange={(e) => console.log('input')}
            type="checkbox"
            checked={state === 'on'}
          />
        </SwitchToggleStyle>
      </div>

      <SwitchFooterStye>
      {automation && (<><div>
           <span css="font-weight: 700">Automation: Enabled</span>
        </div>
        <div>
          <span css="font-weight: 700">Automation Targets: </span>
          <span>
           Min {minValue}{unit} | Max {maxValue}{unit}
          </span>
        </div></>)}
        {schedule && (<><div>
           <span css="font-weight: 700">Schedule: Enabled</span>
        </div>
        <div>
          <span css="font-weight: 700">Schedule Info </span>
          <span>
           every {interval / 60}/h for {duration> 60 ? duration/60 + "h": duration + "min" }
          </span>
        </div></>)}
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
    text-transform: capitalize;
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

async function setSwitch(data: {
  id: string;
  state: 'on' | 'off';
  name: string;
  idx: number;
}) {
  const res = await fetch('/api/switch', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  if (!res.ok) {
    console.log(res);
    throw new Error('err');
  }
  return res.json();
}

export { setSwitch };

export function formatSmartDate(input: Date) {
  const date = input instanceof Date ? input : new Date(input);

  const pad = (n: number) => String(n).padStart(2, "0");

  const now = new Date();

  const isToday =
    date.getDate() === now.getDate() &&
    date.getMonth() === now.getMonth() &&
    date.getFullYear() === now.getFullYear();

  const isThisYear = date.getFullYear() === now.getFullYear();

  const dd = pad(date.getDate());
  const mm = pad(date.getMonth() + 1);
  const yyyy = date.getFullYear();

  const hh = pad(date.getHours());
  const min = pad(date.getMinutes());

  if (isToday) {
    return `${hh}:${min}`;
  }

  if (isThisYear) {
    return `${dd}.${mm} ${hh}:${min}`;
  }

  return `${dd}.${mm}.${yyyy} ${hh}:${min}`;
}

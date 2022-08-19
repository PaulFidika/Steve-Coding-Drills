export default () => {
  const monkey: SampleType = { a: 'DDD' };
  const item: Item = { a: 5, b: 5, c: 10 };

  console.log('hey this works', item);
};

declare global {
  type Monkey = {
    name: string;
  };
}

import { search } from '@/api/Search/search';
import React, { useState } from 'react';

export default function Search() {

  // function debounce(func: (arg0: any) => void | PromiseLike<void>, ms: number | undefined) {
  //   let timeout: string | number | NodeJS.Timeout | undefined;
  //   return (...args: any) =>  {
  //     clearTimeout(timeout);
  //     timeout = setTimeout(() => func(args), ms);
  //   };
  // }
  // const debouncedSetSearch = debounce((value: string) => {Search(value)}, 500)
  let [word, setWord] = useState("")
  let [result, setResult] = useState("")
  const handleSearch = () => {
    const res = search(word)
    setResult(res)
  };


  const handleKeyPress = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === 'Enter') {
      handleSearch();
    }
  };
  return (
    <div >
      <input
        type="text"
        placeholder="Поиск"
        onKeyDown={handleKeyPress}
        onChange={ (event: { target: { value: string } }) => {setWord(event.target.value)}}
      ></input>
      {result ? <div>{result}</div> : <div></div>}
    </div>
  );
}

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
  const handleSearch = () => {
    console.log('Performing search:', searchValue);
    // Здесь логика поиска
    // fetchSearchResults(searchValue);
  };

  const handleInput = (event: React.ChangeEvent<HTMLInputElement>) => {
    setSearchValue(event.target.value);
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
        onChange={ (event: { target: { value: string } }) => {setevent.target.value}}
      ></input>
    </div>
  );
}

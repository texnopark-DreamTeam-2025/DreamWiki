import React, { useState } from 'react';

export default function Search() {

  function debounce(func: (arg0: any) => void | PromiseLike<void>, ms: number | undefined) {
    let timeout: string | number | NodeJS.Timeout | undefined;
    return (...args: any) =>  {
      clearTimeout(timeout);
      timeout = setTimeout(() => func(args), ms);
    };
  }
  const debouncedSetSearch = debounce((value: string) => {Search(value)}, 500)
  let word = useState("")
  return (
    <div className={styles.searchForm}>
      <input
        type="text"
        placeholder="Поиск"
        className={styles.searchForm__input}
        onChange={ (event: { target: { value: string } }) => {debouncedSetSearch(event.target.value)}}
      ></input>
    </div>
  );
}

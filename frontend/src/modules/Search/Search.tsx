import React from 'react';
import search from "Public/static/icons/iconSearch.svg";
import styles from "./Search.module.scss";


export default function Search() {

  // function debounce(func: (arg0: any) => void | PromiseLike<void>, ms: number | undefined) {
  //   let timeout: string | number | NodeJS.Timeout | undefined;
  //   return (...args: any) =>  {
  //     clearTimeout(timeout);
  //     timeout = setTimeout(() => func(args), ms);
  //   };
  // }
  // const debouncedSetSearch = debounce((value: string) => {setSearch(value)}, 500)

  return (
    <div className={styles.searchForm}>
      <img className={styles.searchSearchForm__icon} src={search}></img>
      <input
        type="text"
        placeholder="Поиск"
        className={styles.searchForm__input}
      ></input>
    </div>
  );
}

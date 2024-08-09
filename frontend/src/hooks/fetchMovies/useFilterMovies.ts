import useFetchMovies from './useFetchMovies.ts';
import {getFilterMovies, getSearchMovies, Movies} from '@components/gRPC.tsx';

const useFilterMovies = (initialPage: number, genres: string, yearMin: number, yearMax: number, text: string, path: string) => {
    let limit = 0;
    if (window.innerWidth <= 479) {
        limit = 3;
    } else if (window.innerWidth >= 480 && window.innerWidth <= 1023) {
        limit = 6;
    } else if (window.innerWidth >= 1024 && window.innerWidth <= 1599) {
        limit = 9;
    } else if (window.innerWidth >= 1600 && window.innerWidth <= 2399) {
        limit = 12;
    } else {
        limit = 15;
    }

    const fetchFunctions: Record<string, (page: number) => Promise<Movies[]>> = {
        '/filter': (page: number) => getFilterMovies(0, genres, yearMin, yearMax, limit, page),
        '/search': (page: number) => getSearchMovies(text, limit, page),
    };
    const fetchFunction = fetchFunctions[path] || fetchFunctions['/filter'];

    return useFetchMovies(initialPage, fetchFunction);
};

export default useFilterMovies;

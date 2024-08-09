import useFetchMovies from './useFetchMovies.ts';
import { getMovies } from '@components/gRPC.tsx';

const useMovies = (initialPage: number) => {
    let limit = 0;
    if (window.innerWidth <= 479) {
        limit = 3;
    } else if (window.innerWidth >= 480 && window.innerWidth <= 767) {
        limit = 6;
    } else if (window.innerWidth >= 768 && window.innerWidth <= 1023) {
        limit = 9;
    } else if (window.innerWidth >= 1024 && window.innerWidth <= 1599) {
        limit = 12;
    } else if (window.innerWidth >= 1600 && window.innerWidth <= 2399) {
        limit = 15;
    } else {
        limit = 18;
    }
    const fetchFunction = (page: number) => getMovies(limit, page);

    return useFetchMovies(initialPage, fetchFunction);
};

export default useMovies;
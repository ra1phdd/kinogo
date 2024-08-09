import {getMoviesByTypeMovie, Movies} from '@components/gRPC.tsx';
import useFetchMovies from './useFetchMovies.ts';

const useTypeMovies = (initialPage: number, path: string) => {
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
        '/films': (page: number) => getMoviesByTypeMovie(limit, page, 1),
        '/cartoons': (page: number) => getMoviesByTypeMovie(limit, page, 2),
        '/telecasts': (page: number) => getMoviesByTypeMovie(limit, page, 3),
    };
    const fetchFunction = fetchFunctions[path] || fetchFunctions['/films'];

    return useFetchMovies(initialPage, fetchFunction);
};

export default useTypeMovies;

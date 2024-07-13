import React, {useCallback, useEffect, useState} from "react";
import SearchAside from "@components/specific/Search.tsx";
import FilterAside from "@components/specific/Filter.tsx";
import BestMovieAside from "@components/specific/BestMovie.tsx";
import {Movies, getFilterMovies, getSearchMovies} from '@components/gRPC.tsx';
import {useLocation, useSearchParams} from "react-router-dom";
import MovieCard from "@components/common/MovieCard.tsx";
import '@assets/styles/pages/filter.css';
import {RpcError} from "grpc-web";

// Компонент Home
const Filter: React.FC = () => {
    const [movies, setMovies] = useState<Movies[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(1);
    const location = useLocation();
    const formData = location.state as { genres: string[], year__min: string, year__max: string };
    const [searchParams] = useSearchParams();
    const textForm = searchParams.get('text');

    let text: string;
    if (textForm != null){
        text = textForm
    } else{
        text = ""
    }

    const genres = formData?.genres?.join(', ') || '';
    const yearMin = formData?.year__min || '';
    const yearMax = formData?.year__max || '';

    const fetchFunctions: Record<string, () => Promise<Movies[]>> = {
        '/filter': () => getFilterMovies(genres, Number(yearMin), Number(yearMax), 12, page),
        '/search': () => getSearchMovies(text, 12, page),
    };
    const fetchFunction = fetchFunctions[location.pathname] || fetchFunctions['/filter'];

    const loadMovies = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const newMovies = await fetchFunction();
                setMovies((prevMovies) => [...prevMovies, ...newMovies]);
                setPage((prevPage) => prevPage + 1);
            } catch (error) {
                if (error instanceof RpcError) {
                    if (error.code === 5) {
                        setEnd(true);
                    } else {
                        console.error('Error fetching more movies:', error);
                    }
                } else {
                    console.error('Unexpected error:', error);
                }
            } finally {
                setLoading(false);
            }
        }
    }, [loading, end, fetchFunction]);


    const handleScroll = useCallback(() => {
        if (window.innerHeight + document.documentElement.scrollTop >=
            document.documentElement.offsetHeight - 200) {
            loadMovies();
        }
    }, [loadMovies]);

    useEffect(() => {
        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [handleScroll]);

    useEffect(() => {
        loadMovies();
    }, []);

    return (
        <>
            <div className="sections">
                <div className="section__filters">
                    {location.pathname === "/filter" && (
                        <>
                            <h3>Фильтры: </h3>
                            {genres && <span>{genres}</span>}
                            {yearMax && <span>{"<"}{yearMax}</span>}
                            {yearMin && <span>{">"}{yearMin}</span>}
                        </>
                    )}
                    {location.pathname === "/search" && (
                        <>
                            <h3>Результаты поиска по запросу: </h3>
                            <span>{text}</span>
                        </>
                    )}
                </div>
                <div className="section__videos">
                    {movies.map((movie, index) => (
                        <MovieCard key={movie.id} movie={movie} index={index} limit={12}/>
                    ))}
                    {loading && <p>Загрузка...</p>}
                    {movies.length == 0 && <p>Данных в БД нет</p>}
                </div>
            </div>
            <aside>
                <SearchAside/>
                <FilterAside/>
                <BestMovieAside/>
            </aside>
        </>
    );
}

export default Filter;
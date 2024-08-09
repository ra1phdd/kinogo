import { useCallback, useState } from "react";
import { Movies } from '@components/gRPC.tsx';
import { RpcError } from "grpc-web";

type FetchFunction = (page: number) => Promise<Movies[]>;

const useFetchMovies = (initialPage: number, fetchFunction: FetchFunction) => {
    const [movies, setMovies] = useState<Movies[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(initialPage);

    const loadMovies = useCallback(async () => {
        if (!loading && !end) {
            setLoading(true);
            try {
                const newMovies = await fetchFunction(page);
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
    }, [loading, end, page, fetchFunction]);

    return { movies, loading, end, loadMovies };
};

export default useFetchMovies;

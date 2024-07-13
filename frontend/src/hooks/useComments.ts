import { useState, useEffect, useCallback } from 'react';
import { Comments, getComments } from '@components/gRPC.tsx';

const useComments = (id: number) => {
    const [comments, setComments] = useState<Comments[]>([]);
    const [loading, setLoading] = useState(false);
    const [end, setEnd] = useState(false);
    const [page, setPage] = useState(1);

    const loadMoreComments = useCallback(() => {
        if (!loading && !end) {
            setLoading(true);
            getComments(id, 10, page)
                .then((newComments) => {
                    setComments(prevComments => [...prevComments, ...newComments]);
                    setPage(prevPage => prevPage + 1);
                })
                .catch((error) => {
                    if (error.code === 5) {
                        setEnd(true);
                    } else {
                        console.error('Error fetching comments:', error);
                    }
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [id, loading, page, end]);

    useEffect(() => {
        loadMoreComments();
    }, [loadMoreComments]);

    useEffect(() => {
        const handleScroll = () => {
            if (window.innerHeight + document.documentElement.scrollTop >= document.documentElement.offsetHeight - 100) {
                loadMoreComments();
            }
        };
        window.addEventListener('scroll', handleScroll);
        return () => {
            window.removeEventListener('scroll', handleScroll);
        };
    }, [loadMoreComments]);

    return { comments, loading, loadMoreComments };
};

export default useComments;

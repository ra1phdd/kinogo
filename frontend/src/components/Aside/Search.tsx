import { useEffect, useRef } from "react";
import AnimateElement from "@components/AnimateElement.tsx";

function SearchAside() {
    const asideSearchRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (asideSearchRef.current !== null) {
            AnimateElement(asideSearchRef.current, "animate__fadeInRight", 0);
        }
    }, []);

    return (
        <div className="aside__search" ref={asideSearchRef}>
            <h3>Поиск</h3>
            <form action="/search" method="GET">
                <input type="text" name="text" id="aside__search-input" placeholder="Введите запрос"/>
                <button type="submit"><img src="/src/assets/images/search.svg" alt="Поиск"/></button>
            </form>
            <div className="aside__search-results" id="aside__search-results"></div>
        </div>
    )
}

export default SearchAside
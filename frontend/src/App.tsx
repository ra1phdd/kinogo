import './components/Navigation.tsx'
import {BrowserRouter, Routes, Route} from "react-router-dom";
import './assets/css/src/main.css'
import './assets/css/dist/animate.min.css'
import Navigation from "./components/Navigation.tsx";
import Home from "./pages/Home/Home.tsx"
import Movie from "./pages/Movie/Movie.tsx"

function App() {
    return (
    <>
        <BrowserRouter>
            <div className="container">
                <header className="animate__animated animate__fadeInDown animate__faster" style={{overflow: 'hidden'}}>
                    <div className="header__logotype">
                        <h1><a href="/">KINOTEATR</a></h1>
                    </div>
                    <nav>
                        <Navigation/>
                    </nav>
                </header>
                <main>
                    <Routes>
                        <Route path="/" element={<Home />} />
                        <Route path="/id/:id" element={<Movie />} />
                    </Routes>
                </main>
            </div>
        </BrowserRouter>
    </>
  )
}

export default App

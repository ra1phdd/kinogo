import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import { AuthProvider } from '@/contexts/Auth';
import reportWebVitals from "@/reportWebVitals.ts";

ReactDOM.createRoot(document.getElementById('root')!).render(
  //<React.StrictMode>
      <AuthProvider>
          <App/>
      </AuthProvider>
  //</React.StrictMode>,
)

reportWebVitals(console.log);
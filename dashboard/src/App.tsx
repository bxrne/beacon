import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Navbar } from './components/layout/navbar';
import { HostsPage } from './pages/hosts';
import { TrafficPage } from './pages/traffic';

export default function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-background">
        <Navbar />
        <main className="flex-1">
          <Routes>
            <Route path="/" element={<HostsPage />} />
            <Route path="/traffic" element={<TrafficPage />} />
          </Routes>
        </main>
      </div>
    </BrowserRouter>
  );
}
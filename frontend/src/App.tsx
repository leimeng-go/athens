import { BrowserRouter, Routes, Route } from 'react-router-dom'
import MainLayout from './layouts/MainLayout'
import ModuleListPage from './pages/ModuleListPage'
import UploadPage from './pages/UploadPage'
import StatsPage from './pages/StatsPage'
import SettingsPage from './pages/SettingsPage'

function App() {
  return (
    <BrowserRouter>
       <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<ModuleListPage />} />
          <Route path="upload" element={<UploadPage />} />
          <Route path="stats" element={<StatsPage />} />
          <Route path="settings" element={<SettingsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App

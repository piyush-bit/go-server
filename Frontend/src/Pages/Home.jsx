import React, { useState, useEffect } from 'react';
import { Github, Mail, LogIn, Copy, Check } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';

function Dashboard() {
  const navigate = useNavigate();
  const location = useLocation();
  const [homeData, setHomeData] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingApp, setEditingApp] = useState(null);
  const [formData, setFormData] = useState({ name: '', callback_url: '' });
  const [alert, setAlert] = useState({ show: false, message: '', type: 'success' });
  const [copiedId, setCopiedId] = useState(null);

  const [refresh, setRefresh] = useState(0);

  const showAlert = (message, type = 'success') => {
    setAlert({ show: true, message, type });
    setTimeout(() => setAlert({ show: false, message: '', type: 'success' }), 3000);
  };

  const BACKEND_URI = import.meta.env.VITE_BACKEND_URI ?? "";

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch(BACKEND_URI + '/api/v1/app/', {
          headers: {
            'Authorization': 'Bearer ' + localStorage.getItem('token')
          }
        });
        
        if (!response.ok) {
          if (response.status === 401) {
            navigate('/')
          }
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        setHomeData(data.data);
      } catch (error) {
        console.error('Error fetching data:', error);
        showAlert('Failed to load applications', 'error');
      }
    };

    fetchData();
  }, [refresh]);

  const handleDelete = async (id) => {
    try {
      const response = await fetch(BACKEND_URI + `/api/v1/app/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        }
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setHomeData(prev => ({
        ...prev,
        apps: prev.apps.filter(app => app.id !== id)
      }));
      showAlert('Application deleted successfully');
    } catch (error) {
      console.error('Error deleting app:', error);
      showAlert('Failed to delete application', 'error');
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const url = editingApp 
        ? BACKEND_URI + `/api/v1/app/${editingApp.id}`
        : BACKEND_URI + '/api/v1/app/create';
      
      const formDataToSend = new FormData();
      formDataToSend.append('name', formData.name);
      formDataToSend.append('callback_url', formData.callback_url);

      const response = await fetch(url, {
        method: editingApp ? 'PATCH' : 'POST',
        headers: {
          'Authorization': 'Bearer ' + localStorage.getItem('token')
        },
        body: formDataToSend
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      
      if (editingApp) {
        setHomeData(prev => ({
          ...prev,
          apps: prev.apps.map(app => app.id === editingApp.id ? data.data : app)
        }));
        showAlert('Application updated successfully');
      } else {
        setHomeData(prev => ({
          ...prev,
          apps: [...prev.apps, data.data]
        }));
        showAlert('Application created successfully');
      }
      
      setIsModalOpen(false);
      setEditingApp(null);
      setFormData({ name: '', callback_url: '' });
    } catch (error) {
      console.error('Error submitting form:', error);
      showAlert(`Failed to ${editingApp ? 'update' : 'create'} application`, 'error');
    }
  };

  const openEditModal = (app) => {
    setEditingApp(app);
    setFormData({ name: app.name, callback_url: app.callback_url });
    setIsModalOpen(true);
  };

  const copyUrl = (id) => {
    const baseUrl = window.location.origin;
    const urlToCopy = `${baseUrl}?id=${id}`;
    
    navigator.clipboard.writeText(urlToCopy)
      .then(() => {
        setCopiedId(id);
        setTimeout(() => setCopiedId(null), 2000);
        showAlert('URL copied to clipboard');
      })
      .catch(err => {
        console.error('Failed to copy: ', err);
        showAlert('Failed to copy URL', 'error');
      });
  };

  if (!homeData) return <div className="flex items-center justify-center h-screen">Loading...</div>;

  return (
    <div className="p-6 max-w-7xl w-full mx-auto">
      {alert.show && (
        <div 
          className={`mb-4 p-4 rounded-lg ${
            alert.type === 'error' 
              ? 'bg-red-100 border border-red-400 text-red-700' 
              : 'bg-green-100 border border-green-400 text-green-700'
          }`}
        >
          {alert.message}
        </div>
      )}

      <div className="flex justify-between items-center mb-8">
        <h1 className="text-2xl font-bold">SSO Dashboard</h1>
        <button 
          onClick={() => {
            setEditingApp(null);
            setFormData({ name: '', callback_url: '' });
            setIsModalOpen(true);
          }}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition duration-200"
        >
          Add New App
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
        {homeData.apps?.map(app => (
          <div key={app.id} className="bg-white rounded-lg shadow p-6 border w-full border-gray-200">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-semibold">{app.name}</h3>
              <div className="flex gap-2">
                <button
                  onClick={() => openEditModal(app)}
                  className="px-3 py-1 bg-gray-600 text-white rounded hover:bg-gray-700 transition duration-200"
                >
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(app.id)}
                  className="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700 transition duration-200"
                >
                  Delete
                </button>
              </div>
            </div>
            <p className="text-gray-600 break-all mb-4">{app.callback_url}</p>
            <div className="flex items-center">
              <button
                onClick={() => copyUrl(app.id)}
                className="flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition duration-200"
              >
                {copiedId === app.id ? (
                  <>
                    <Check size={16} />
                    Copied!
                  </>
                ) : (
                  <>
                    <Copy size={16} />
                    Copy URL
                  </>
                )}
              </button>
            </div>
          </div>
        ))}
      </div>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold">
                {editingApp ? 'Edit Application' : 'New Application'}
              </h2>
              <button 
                onClick={() => setIsModalOpen(false)}
                className="text-gray-500 hover:text-gray-700"
              >
                ×
              </button>
            </div>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Application Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter application name"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Callback URL
                </label>
                <input
                  type="text"
                  value={formData.callback_url}
                  onChange={(e) => setFormData(prev => ({ ...prev, callback_url: e.target.value }))}
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter callback URL"
                  required
                />
              </div>
              <div className="flex justify-end gap-2">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300 transition duration-200"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition duration-200"
                >
                  {editingApp ? 'Update' : 'Create'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Dashboard;
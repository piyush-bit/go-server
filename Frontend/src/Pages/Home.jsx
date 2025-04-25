import React, { useState, useEffect } from 'react';
import { Github, Mail, LogIn, Copy, Check, Plus, Edit, Trash2, X, Loader } from 'lucide-react';
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
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [refresh, setRefresh] = useState(0);

  const showAlert = (message, type = 'success') => {
    setAlert({ show: true, message, type });
    setTimeout(() => setAlert({ show: false, message: '', type: 'success' }), 3000);
  };

  const BACKEND_URI = import.meta.env.VITE_BACKEND_URI ?? "";

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
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
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, [refresh, navigate, BACKEND_URI]);

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
    setIsSubmitting(true);
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
          apps: [...(prev.apps??{}), data.data]
        }));
        showAlert('Application created successfully');
      }
      
      setIsModalOpen(false);
      setEditingApp(null);
      setFormData({ name: '', callback_url: '' });
    } catch (error) {
      console.error('Error submitting form:', error);
      showAlert(`Failed to ${editingApp ? 'update' : 'create'} application`, 'error');
    } finally {
      setIsSubmitting(false);
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

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        {/* Alert Component */}
        {alert.show && (
          <div className="fixed top-4 right-4 z-50 w-72 transition-all duration-300 ease-in-out transform translate-y-0 opacity-100">
            <div 
              className={`p-4 rounded-md shadow-lg border-l-4 flex items-center ${
                alert.type === 'error' 
                  ? 'bg-red-50 border-red-500 text-red-800' 
                  : 'bg-green-50 border-green-500 text-green-800'
              }`}
            >
              <div className="flex-shrink-0 mr-3">
                {alert.type === 'error' ? (
                  <X className="h-5 w-5 text-red-500" />
                ) : (
                  <Check className="h-5 w-5 text-green-500" />
                )}
              </div>
              <div className="text-sm font-medium">{alert.message}</div>
            </div>
          </div>
        )}

        {/* Header */}
        <div className="mb-10">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-extrabold text-gray-900">SSO Dashboard</h1>
            <button 
              onClick={() => {
                setEditingApp(null);
                setFormData({ name: '', callback_url: '' });
                setIsModalOpen(true);
              }}
              className="inline-flex items-center px-4 py-2 bg-blue-600 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-200"
            >
              <Plus className="h-4 w-4 mr-2" />
              Add New App
            </button>
          </div>
          <p className="mt-2 text-sm text-gray-600">Manage your single sign-on applications</p>
        </div>

        {/* Content */}
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="flex flex-col items-center">
              <Loader className="h-8 w-8 text-blue-500 animate-spin" />
              <p className="mt-4 text-gray-600">Loading applications...</p>
            </div>
          </div>
        ) : homeData?.apps?.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {homeData.apps.map(app => (
              <div key={app.id} className="bg-white rounded-lg shadow overflow-hidden border border-gray-200 hover:shadow-md transition-shadow duration-300">
                <div className="px-6 py-5 border-b border-gray-200 bg-gray-50">
                  <div className="flex justify-between items-center">
                    <h3 className="text-lg font-medium text-gray-900">{app.name}</h3>
                    <div className="flex space-x-2">
                      <button
                        onClick={() => openEditModal(app)}
                        className="p-1 rounded-full hover:bg-gray-200 text-gray-600 transition-colors duration-200"
                        title="Edit application"
                      >
                        <Edit className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => handleDelete(app.id)}
                        className="p-1 rounded-full hover:bg-red-100 text-red-600 transition-colors duration-200"
                        title="Delete application"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                </div>
                <div className="px-6 py-4">
                  <div className="text-sm font-medium text-gray-600 mb-1">Callback URL</div>
                  <p className="text-sm text-gray-800 break-all mb-5 bg-gray-50 p-3 rounded border border-gray-200">{app.callback_url}</p>
                  <button
                    onClick={() => copyUrl(app.id)}
                    className="inline-flex items-center px-3 py-2 text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 transition-colors duration-200 w-full justify-center"
                  >
                    {copiedId === app.id ? (
                      <>
                        <Check className="h-4 w-4 mr-2" />
                        URL Copied
                      </>
                    ) : (
                      <>
                        <Copy className="h-4 w-4 mr-2" />
                        Copy Login URL
                      </>
                    )}
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-64 bg-white rounded-lg shadow p-8 border border-gray-200">
            <div className="text-center">
              <h3 className="text-lg font-medium text-gray-900 mb-2">No applications found</h3>
              <p className="text-gray-600 mb-6">Get started by adding your first application</p>
              <button 
                onClick={() => {
                  setEditingApp(null);
                  setFormData({ name: '', callback_url: '' });
                  setIsModalOpen(true);
                }}
                className="inline-flex items-center px-4 py-2 bg-blue-600 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-200"
              >
                <Plus className="h-4 w-4 mr-2" />
                Add New App
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-75 flex items-center justify-center p-4 z-50">
          <div 
            className="bg-white rounded-lg shadow-xl w-full max-w-md transform transition-all"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="px-6 py-4 border-b border-gray-200">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-medium text-gray-900">
                  {editingApp ? 'Edit Application' : 'Create New Application'}
                </h2>
                <button 
                  onClick={() => setIsModalOpen(false)}
                  className="text-gray-400 hover:text-gray-600 focus:outline-none"
                >
                  <X className="h-5 w-5" />
                </button>
              </div>
            </div>
            <form onSubmit={handleSubmit}>
              <div className="px-6 py-4 space-y-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Application Name
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
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
                    className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="https://your-app.com/callback"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">
                    The URL where users will be redirected after authentication
                  </p>
                </div>
              </div>
              <div className="px-6 py-4 bg-gray-50 border-t border-gray-200 flex justify-end space-x-3 rounded-b-lg">
                <button
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  disabled={isSubmitting}
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? (
                    <>
                      <Loader className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" />
                      {editingApp ? 'Updating...' : 'Creating...'}
                    </>
                  ) : (
                    <>{editingApp ? 'Update Application' : 'Create Application'}</>
                  )}
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
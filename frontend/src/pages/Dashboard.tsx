import React, { useState, useEffect } from 'react';
import API from '../services/api';
import { 
  LayoutDashboard, 
  Users, 
  GraduationCap, 
  BookOpen, 
  Award, 
  LogOut, 
  Search, 
  Bell, 
  TrendingUp, 
  ShieldCheck, 
  Database,
  Activity,
  Plus,
  X,
  CheckCircle2
} from 'lucide-react';

export default function Dashboard() {
  const [activeTab, setActiveTab] = useState('overview');
  const [results, setResults] = useState([]);
  const [students, setStudents] = useState([]);
  const [teachers, setTeachers] = useState([]);
  
  const [searchSeat, setSearchSeat] = useState('');
  const [searchResult, setSearchResult] = useState<any>(null);
  const [loading, setLoading] = useState(false);
  const [successMsg, setSuccessMsg] = useState('');

  // Modals state
  const [showResultModal, setShowResultModal] = useState(false);
  const [showStudentModal, setShowStudentModal] = useState(false);
  const [showTeacherModal, setShowTeacherModal] = useState(false);

  // Form states
  const [resultForm, setResultForm] = useState({
    seat_number: '',
    student_name: '',
    section: 'علمي',
    governorate: 'صنعاء',
    total_score: '',
    percentage: ''
  });

  const [studentForm, setStudentForm] = useState({
    name: '',
    email: '',
    grade: 'الثالث الثانوي',
    section: 'علمي'
  });

  const [teacherForm, setTeacherForm] = useState({
    name: '',
    email: '',
    subject: 'الرياضيات'
  });

  useEffect(() => {
    if (activeTab === 'results') {
      fetchResults();
    } else if (activeTab === 'students') {
      fetchStudents();
    } else if (activeTab === 'teachers') {
      fetchTeachers();
    }
  }, [activeTab]);

  const fetchResults = async () => {
    try {
      const res = await API.get('/results');
      setResults(res.data.data || []);
    } catch (err) {
      console.error('Failed to fetch results', err);
    }
  };

  const fetchStudents = async () => {
    try {
      const res = await API.get('/students');
      setStudents(res.data.data || []);
    } catch (err) {
      console.error('Failed to fetch students', err);
    }
  };

  const fetchTeachers = async () => {
    try {
      const res = await API.get('/teachers');
      setTeachers(res.data.data || []);
    } catch (err) {
      console.error('Failed to fetch teachers', err);
    }
  };

  const handleSearchBySeat = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchSeat) return;
    setLoading(true);
    try {
      const res = await API.get(`/results/seat/${searchSeat}`);
      setSearchResult(res.data.data);
    } catch (err) {
      setSearchResult(null);
      alert('لم يتم العثور على نتيجة برقم الجلوس هذا');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateResult = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await API.post('/results', {
        ...resultForm,
        total_score: parseFloat(resultForm.total_score),
        percentage: parseFloat(resultForm.percentage)
      });
      setSuccessMsg('تمت إضافة النتيجة بنجاح!');
      setShowResultModal(false);
      setResultForm({ seat_number: '', student_name: '', section: 'علمي', governorate: 'صنعاء', total_score: '', percentage: '' });
      fetchResults();
      setTimeout(() => setSuccessMsg(''), 4000);
    } catch (err) {
      alert('فشل إضافة النتيجة، تأكد من البيانات');
    }
  };

  const handleCreateStudent = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await API.post('/students', studentForm);
      setSuccessMsg('تمت إضافة الطالب بنجاح!');
      setShowStudentModal(false);
      setStudentForm({ name: '', email: '', grade: 'الثالث الثانوي', section: 'علمي' });
      fetchStudents();
      setTimeout(() => setSuccessMsg(''), 4000);
    } catch (err) {
      alert('فشل إضافة الطالب');
    }
  };

  const handleCreateTeacher = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await API.post('/teachers', teacherForm);
      setSuccessMsg('تمت إضافة المعلم بنجاح!');
      setShowTeacherModal(false);
      setTeacherForm({ name: '', email: '', subject: 'الرياضيات' });
      fetchTeachers();
      setTimeout(() => setSuccessMsg(''), 4000);
    } catch (err) {
      alert('فشل إضافة المعلم');
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    window.location.reload();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-blue-950 text-slate-100 flex font-sans dir=rtl">
      
      {/* 📊 Sidebar الفخم */}
      <aside className="w-72 bg-slate-900/90 backdrop-blur-2xl border-l border-slate-800 p-6 flex flex-col justify-between shadow-2xl z-20">
        <div>
          <div className="flex items-center gap-3 mb-10 px-2">
            <div className="bg-gradient-to-tr from-blue-600 to-indigo-500 p-3 rounded-2xl shadow-lg shadow-blue-500/30 text-white">
              <GraduationCap size={28} />
            </div>
            <div>
              <h1 className="font-bold text-lg text-white tracking-wide">وزارة التعليم</h1>
              <p className="text-xs text-blue-400 font-medium">نظام الإدارة المركزية</p>
            </div>
          </div>

          <nav className="space-y-2">
            <button
              onClick={() => setActiveTab('overview')}
              className={`w-full flex items-center gap-3.5 px-4 py-3 rounded-xl font-medium transition-all duration-300 ${
                activeTab === 'overview'
                  ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30'
                  : 'text-slate-400 hover:bg-slate-800/60 hover:text-white'
              }`}
            >
              <LayoutDashboard size={20} />
              <span>نظرة عامة</span>
            </button>

            <button
              onClick={() => setActiveTab('results')}
              className={`w-full flex items-center gap-3.5 px-4 py-3 rounded-xl font-medium transition-all duration-300 ${
                activeTab === 'results'
                  ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30'
                  : 'text-slate-400 hover:bg-slate-800/60 hover:text-white'
              }`}
            >
              <Award size={20} />
              <span>نتائج الثانوية العامة</span>
            </button>

            <button
              onClick={() => setActiveTab('students')}
              className={`w-full flex items-center gap-3.5 px-4 py-3 rounded-xl font-medium transition-all duration-300 ${
                activeTab === 'students'
                  ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30'
                  : 'text-slate-400 hover:bg-slate-800/60 hover:text-white'
              }`}
            >
              <Users size={20} />
              <span>إدارة الطلاب</span>
            </button>

            <button
              onClick={() => setActiveTab('teachers')}
              className={`w-full flex items-center gap-3.5 px-4 py-3 rounded-xl font-medium transition-all duration-300 ${
                activeTab === 'teachers'
                  ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30'
                  : 'text-slate-400 hover:bg-slate-800/60 hover:text-white'
              }`}
            >
              <BookOpen size={20} />
              <span>إدارة المعلمين</span>
            </button>
          </nav>
        </div>

        <div className="pt-6 border-t border-slate-800/80">
          <button
            onClick={handleLogout}
            className="w-full flex items-center justify-center gap-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/20 py-2.5 rounded-xl transition-all font-medium text-sm"
          >
            <LogOut size={18} />
            <span>تسجيل الخروج</span>
          </button>
        </div>
      </aside>

      {/* 🚀 Main Content Area */}
      <main className="flex-1 flex flex-col overflow-y-auto">
        
        {/* Top Navbar */}
        <header className="h-20 bg-slate-900/40 backdrop-blur-md border-b border-slate-800/80 px-8 flex items-center justify-between sticky top-0 z-10">
          <div className="flex items-center gap-3">
            <span className="bg-emerald-500/10 text-emerald-400 text-xs font-semibold px-3 py-1.5 rounded-full border border-emerald-500/20 flex items-center gap-2 shadow-inner">
              <span className="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse"></span>
              النظام متصل وآمن (Redis & JWT)
            </span>
            {successMsg && (
              <span className="bg-blue-500/10 text-blue-400 text-xs font-semibold px-3 py-1.5 rounded-full border border-blue-500/20 flex items-center gap-1.5 animate-bounce">
                <CheckCircle2 size={14} /> {successMsg}
              </span>
            )}
          </div>

          <div className="flex items-center gap-4">
            <button className="p-2.5 bg-slate-800/60 hover:bg-slate-800 rounded-xl text-slate-300 transition border border-slate-700/50 relative">
              <Bell size={20} />
              <span className="absolute top-2 left-2 w-2 h-2 bg-blue-500 rounded-full"></span>
            </button>
            <div className="flex items-center gap-3 bg-slate-800/50 border border-slate-700/50 px-4 py-2 rounded-xl">
              <div className="w-9 h-9 rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center font-bold text-white shadow-md">
                أ
              </div>
              <div>
                <h4 className="text-sm font-semibold text-white">مدير النظام</h4>
                <p className="text-xs text-slate-400">admin@moe.gov</p>
              </div>
            </div>
          </div>
        </header>

        {/* Dashboard Body */}
        <div className="p-8 max-w-7xl mx-auto w-full">
          
          {/* OVERVIEW TAB */}
          {activeTab === 'overview' && (
            <div className="space-y-8 animate-fadeIn">
              <div className="relative overflow-hidden bg-gradient-to-r from-blue-600 via-indigo-600 to-violet-700 rounded-3xl p-8 shadow-2xl shadow-blue-500/20 border border-blue-400/20">
                <div className="absolute -left-10 -bottom-10 w-64 h-64 bg-white/10 rounded-full blur-3xl"></div>
                <div className="relative z-10">
                  <h2 className="text-3xl font-extrabold text-white mb-2">أهلاً بك، المدير العام 👋</h2>
                  <p className="text-blue-100 max-w-xl text-sm leading-relaxed">
                    لوحة تحكم نظام وزارة التعليم المركزي. يمكنك إدارة الطلاب، المعلمين، وإدخال واستعراض نتائج الثانوية العامة بمرونة تامة.
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                  <div className="flex items-center justify-between mb-4">
                    <div className="p-3 bg-blue-500/10 text-blue-400 rounded-xl"><Users size={24} /></div>
                    <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-lg font-medium flex items-center gap-1"><TrendingUp size={14} /> نشط</span>
                  </div>
                  <h3 className="text-slate-400 text-sm font-medium mb-1">إجمالي الطلاب</h3>
                  <p className="text-3xl font-bold text-white">1,240</p>
                </div>

                <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                  <div className="flex items-center justify-between mb-4">
                    <div className="p-3 bg-indigo-500/10 text-indigo-400 rounded-xl"><BookOpen size={24} /></div>
                    <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-lg font-medium flex items-center gap-1"><TrendingUp size={14} /> نشط</span>
                  </div>
                  <h3 className="text-slate-400 text-sm font-medium mb-1">المعلمون النشطون</h3>
                  <p className="text-3xl font-bold text-white">85</p>
                </div>

                <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                  <div className="flex items-center justify-between mb-4">
                    <div className="p-3 bg-emerald-500/10 text-emerald-400 rounded-xl"><Award size={24} /></div>
                    <span className="text-xs text-emerald-400 bg-emerald-500/10 px-2.5 py-1 rounded-lg font-medium flex items-center gap-1"><Activity size={14} /> محدث</span>
                  </div>
                  <h3 className="text-slate-400 text-sm font-medium mb-1">نتائج الثانوية العامة</h3>
                  <p className="text-3xl font-bold text-white">5,420</p>
                </div>

                <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                  <div className="flex items-center justify-between mb-4">
                    <div className="p-3 bg-violet-500/10 text-violet-400 rounded-xl"><ShieldCheck size={24} /></div>
                    <span className="text-xs text-blue-400 bg-blue-500/10 px-2.5 py-1 rounded-lg font-medium">RBAC مفعل</span>
                  </div>
                  <h3 className="text-slate-400 text-sm font-medium mb-1">مستوى الأمان</h3>
                  <p className="text-3xl font-bold text-white">100%</p>
                </div>
              </div>
            </div>
          )}

          {/* RESULTS TAB */}
          {activeTab === 'results' && (
            <div className="space-y-6 animate-fadeIn">
              <div className="flex flex-col md:flex-row items-center justify-between gap-4 bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                <div>
                  <h3 className="text-xl font-bold text-white">إدارة نتائج الثانوية العامة</h3>
                  <p className="text-slate-400 text-sm">البحث السريع برقم الجلوس أو إضافة نتيجة جديدة لقاعدة البيانات</p>
                </div>
                
                <div className="flex items-center gap-3 w-full md:w-auto">
                  <button
                    onClick={() => setShowResultModal(true)}
                    className="bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white px-5 py-2.5 rounded-xl font-medium transition shadow-lg shadow-blue-600/20 text-sm flex items-center gap-2"
                  >
                    <Plus size={18} /> إضافة نتيجة جديدة
                  </button>
                </div>
              </div>

              {/* Search Bar */}
              <form onSubmit={handleSearchBySeat} className="flex gap-3 bg-slate-900/60 p-4 rounded-2xl border border-slate-800">
                <div className="relative flex-1">
                  <span className="absolute inset-y-0 right-0 pr-3 flex items-center text-slate-400"><Search size={18} /></span>
                  <input
                    type="text"
                    placeholder="بحث برقم الجلوس (مثال: 901122)"
                    value={searchSeat}
                    onChange={(e) => setSearchSeat(e.target.value)}
                    className="w-full pr-10 pl-4 py-2 bg-slate-800/80 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                  />
                </div>
                <button type="submit" className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-xl text-sm font-medium transition">بحث</button>
              </form>

              {searchResult && (
                <div className="bg-gradient-to-br from-slate-900 to-blue-950 border border-blue-500/30 p-6 rounded-2xl shadow-2xl">
                  <h4 className="text-lg font-bold text-blue-400 mb-4 flex items-center gap-2">
                    <Award size={20} /> نتيجة الطالب: {searchResult.student_name}
                  </h4>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="bg-slate-800/50 p-4 rounded-xl border border-slate-700/50">
                      <span className="text-slate-400 text-xs block mb-1">رقم الجلوس</span>
                      <span className="text-white font-bold">{searchResult.seat_number}</span>
                    </div>
                    <div className="bg-slate-800/50 p-4 rounded-xl border border-slate-700/50">
                      <span className="text-slate-400 text-xs block mb-1">المجموع الكلي</span>
                      <span className="text-emerald-400 font-bold">{searchResult.total_score}</span>
                    </div>
                    <div className="bg-slate-800/50 p-4 rounded-xl border border-slate-700/50">
                      <span className="text-slate-400 text-xs block mb-1">النسبة المئوية</span>
                      <span className="text-blue-400 font-bold">{searchResult.percentage}%</span>
                    </div>
                    <div className="bg-slate-800/50 p-4 rounded-xl border border-slate-700/50">
                      <span className="text-slate-400 text-xs block mb-1">المحافظة / الشعبة</span>
                      <span className="text-white font-bold">{searchResult.governorate} - {searchResult.section}</span>
                    </div>
                  </div>
                </div>
              )}

              {/* Results Table */}
              <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 rounded-2xl shadow-xl overflow-hidden">
                <div className="p-6 border-b border-slate-800">
                  <h4 className="font-bold text-white text-lg">سجلات النتائج</h4>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-right">
                    <thead>
                      <tr className="bg-slate-800/50 text-slate-400 text-xs uppercase tracking-wider">
                        <th className="py-4 px-6">رقم الجلوس</th>
                        <th className="py-4 px-6">اسم الطالب</th>
                        <th className="py-4 px-6">الشعبة</th>
                        <th className="py-4 px-6">المحافظة</th>
                        <th className="py-4 px-6">المجموع</th>
                        <th className="py-4 px-6">النسبة</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800 text-sm">
                      {results.length > 0 ? (
                        results.map((res: any) => (
                          <tr key={res.id} className="hover:bg-slate-800/30 transition">
                            <td className="py-4 px-6 font-mono text-blue-400">{res.seat_number}</td>
                            <td className="py-4 px-6 font-semibold text-white">{res.student_name}</td>
                            <td className="py-4 px-6 text-slate-300">{res.section}</td>
                            <td className="py-4 px-6 text-slate-300">{res.governorate}</td>
                            <td className="py-4 px-6 font-bold text-emerald-400">{res.total_score}</td>
                            <td className="py-4 px-6 font-bold text-blue-400">{res.percentage}%</td>
                          </tr>
                        ))
                      ) : (
                        <tr><td colSpan={6} className="text-center py-8 text-slate-500">لا توجد نتائج مسجلة حتى الآن.</td></tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* STUDENTS TAB */}
          {activeTab === 'students' && (
            <div className="space-y-6 animate-fadeIn">
              <div className="flex items-center justify-between bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                <div>
                  <h3 className="text-xl font-bold text-white">إدارة الطلاب</h3>
                  <p className="text-slate-400 text-sm">استعراض وإضافة الطلاب المسجلين في المنظومة التعليمية</p>
                </div>
                <button
                  onClick={() => setShowStudentModal(true)}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl font-medium transition shadow-lg text-sm flex items-center gap-2"
                >
                  <Plus size={18} /> إضافة طالب جديد
                </button>
              </div>

              <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 rounded-2xl shadow-xl overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full text-right">
                    <thead>
                      <tr className="bg-slate-800/50 text-slate-400 text-xs uppercase">
                        <th className="py-4 px-6">اسم الطالب</th>
                        <th className="py-4 px-6">البريد الإلكتروني</th>
                        <th className="py-4 px-6">المرحلة الدراسية</th>
                        <th className="py-4 px-6">الشعبة</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800 text-sm">
                      {students.length > 0 ? (
                        students.map((st: any) => (
                          <tr key={st.id} className="hover:bg-slate-800/30 transition">
                            <td className="py-4 px-6 font-semibold text-white">{st.name}</td>
                            <td className="py-4 px-6 text-blue-400">{st.email}</td>
                            <td className="py-4 px-6 text-slate-300">{st.grade}</td>
                            <td className="py-4 px-6 text-slate-300">{st.section}</td>
                          </tr>
                        ))
                      ) : (
                        <tr><td colSpan={4} className="text-center py-8 text-slate-500">لا يوجد طلاب مسجلون حالياً.</td></tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* TEACHERS TAB */}
          {activeTab === 'teachers' && (
            <div className="space-y-6 animate-fadeIn">
              <div className="flex items-center justify-between bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 p-6 rounded-2xl shadow-xl">
                <div>
                  <h3 className="text-xl font-bold text-white">إدارة المعلمين</h3>
                  <p className="text-slate-400 text-sm">إدارة الطاقم التعليمي وتخصصاتهم الأكاديمية</p>
                </div>
                <button
                  onClick={() => setShowTeacherModal(true)}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-xl font-medium transition shadow-lg text-sm flex items-center gap-2"
                >
                  <Plus size={18} /> إضافة معلم جديد
                </button>
              </div>

              <div className="bg-slate-900/60 backdrop-blur-xl border border-slate-800/80 rounded-2xl shadow-xl overflow-hidden">
                <div className="overflow-x-auto">
                  <table className="w-full text-right">
                    <thead>
                      <tr className="bg-slate-800/50 text-slate-400 text-xs uppercase">
                        <th className="py-4 px-6">اسم المعلم</th>
                        <th className="py-4 px-6">البريد الإلكتروني</th>
                        <th className="py-4 px-6">المادة الدراسية</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800 text-sm">
                      {teachers.length > 0 ? (
                        teachers.map((tc: any) => (
                          <tr key={tc.id} className="hover:bg-slate-800/30 transition">
                            <td className="py-4 px-6 font-semibold text-white">{tc.name}</td>
                            <td className="py-4 px-6 text-blue-400">{tc.email}</td>
                            <td className="py-4 px-6 text-emerald-400 font-medium">{tc.subject}</td>
                          </tr>
                        ))
                      ) : (
                        <tr><td colSpan={3} className="text-center py-8 text-slate-500">لا يوجد معلمون مسجلون حالياً.</td></tr>
                      )}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

        </div>
      </main>

      {/* --- MODALS --- */}

      {/* Result Modal */}
      {showResultModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl p-8 max-w-lg w-full shadow-2xl relative">
            <button onClick={() => setShowResultModal(false)} className="absolute top-6 left-6 text-slate-400 hover:text-white"><X size={20} /></button>
            <h3 className="text-xl font-bold text-white mb-4">إضافة نتيجة ثانوية عامة جديدة</h3>
            <form onSubmit={handleCreateResult} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">رقم الجلوس</label>
                <input type="text" required value={resultForm.seat_number} onChange={e => setResultForm({...resultForm, seat_number: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" placeholder="مثال: 901133" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">اسم الطالب الكامل</label>
                <input type="text" required value={resultForm.student_name} onChange={e => setResultForm({...resultForm, student_name: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" placeholder="مثال: محمد أحمد علي" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">الشعبة</label>
                  <select value={resultForm.section} onChange={e => setResultForm({...resultForm, section: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm">
                    <option value="علمي">علمي</option>
                    <option value="أدبي">أدبي</option>
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">المحافظة</label>
                  <input type="text" required value={resultForm.governorate} onChange={e => setResultForm({...resultForm, governorate: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">المجموع الكلي</label>
                  <input type="number" step="any" required value={resultForm.total_score} onChange={e => setResultForm({...resultForm, total_score: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" placeholder="مثال: 485" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">النسبة المئوية %</label>
                  <input type="number" step="any" required value={resultForm.percentage} onChange={e => setResultForm({...resultForm, percentage: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" placeholder="مثال: 97.0" />
                </div>
              </div>
              <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-xl font-semibold transition mt-4 shadow-lg shadow-blue-600/30">حفظ النتيجة في قاعدة البيانات</button>
            </form>
          </div>
        </div>
      )}

      {/* Student Modal */}
      {showStudentModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl p-8 max-w-lg w-full shadow-2xl relative">
            <button onClick={() => setShowStudentModal(false)} className="absolute top-6 left-6 text-slate-400 hover:text-white"><X size={20} /></button>
            <h3 className="text-xl font-bold text-white mb-4">إضافة طالب جديد</h3>
            <form onSubmit={handleCreateStudent} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">اسم الطالب</label>
                <input type="text" required value={studentForm.name} onChange={e => setStudentForm({...studentForm, name: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">البريد الإلكتروني</label>
                <input type="email" required value={studentForm.email} onChange={e => setStudentForm({...studentForm, email: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">المرحلة</label>
                  <input type="text" required value={studentForm.grade} onChange={e => setStudentForm({...studentForm, grade: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-400 mb-1">الشعبة</label>
                  <input type="text" required value={studentForm.section} onChange={e => setStudentForm({...studentForm, section: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
                </div>
              </div>
              <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-xl font-semibold transition mt-4 shadow-lg">حفظ الطالب</button>
            </form>
          </div>
        </div>
      )}

      {/* Teacher Modal */}
      {showTeacherModal && (
        <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
          <div className="bg-slate-900 border border-slate-800 rounded-3xl p-8 max-w-lg w-full shadow-2xl relative">
            <button onClick={() => setShowTeacherModal(false)} className="absolute top-6 left-6 text-slate-400 hover:text-white"><X size={20} /></button>
            <h3 className="text-xl font-bold text-white mb-4">إضافة معلم جديد</h3>
            <form onSubmit={handleCreateTeacher} className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">اسم المعلم</label>
                <input type="text" required value={teacherForm.name} onChange={e => setTeacherForm({...teacherForm, name: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">البريد الإلكتروني</label>
                <input type="email" required value={teacherForm.email} onChange={e => setTeacherForm({...teacherForm, email: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
              </div>
              <div>
                <label className="block text-xs font-semibold text-slate-400 mb-1">المادة الدراسية</label>
                <input type="text" required value={teacherForm.subject} onChange={e => setTeacherForm({...teacherForm, subject: e.target.value})} className="w-full px-4 py-2.5 bg-slate-800 border border-slate-700 rounded-xl text-white text-sm" />
              </div>
              <button type="submit" className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-xl font-semibold transition mt-4 shadow-lg">حفظ المعلم</button>
            </form>
          </div>
        </div>
      )}

    </div>
  );
}

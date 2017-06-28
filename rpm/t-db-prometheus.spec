Name: t-db-prometheus
Version:1.0.1
Release: %(echo $RELEASE)
# if you want use the parameter of rpm_create on build time,
# uncomment below
Summary: Alibaba DB Monitor Server.
Group: alibaba/application
License: copyright by alibaba inc.
%define _prefix /opt/prometheus
%define _unitdir /usr/lib/systemd/system
%define shortname prometheus
%define debug_package %{nil}


BuildRequires: t-dbfree-golang

# uncomment below, if depend on other packages

#Requires: package_name = 1.0.0

%description
# if you want publish current svn URL or Revision use these macros
Alibaba DB Syreo Monitor Service.

# %debug_package
# support debuginfo package, to reduce runtime package size

# prepare your files
%install
# OLDPWD is the dir of rpm_create running
# _prefix is an inner var of rpmbuild,
# can set by rpm_create, default is "/home/a"
# _lib is an inner var, maybe "lib" or "lib64" depend on OS

# create dirs
mkdir -p $RPM_BUILD_ROOT%{_prefix}
mkdir -p $RPM_BUILD_ROOT%{_prefix}/bin


#install
cd -
cp ../../src/github.com/qinguoan/prometheus/cmd/prometheus/prometheus $RPM_BUILD_ROOT%{_prefix}/bin/

rm -rf %{buildroot}%{_sysconfdir}/%{shortname}/prometheus
rm -rf %{buildroot}%{_unitdir}/prometheus.service
rm -rf %{buildroot}%{_sysconfdir}/%{shortname}/rules

install -d -m 0755 %{buildroot}%{_sysconfdir}/%{shortname}
install -m 644 -t %{buildroot}%{_sysconfdir}/%{shortname} systemd/environ/prometheus
install -m 644 -t %{buildroot}%{_sysconfdir}/%{shortname} systemd/environ/prometheus.yml
cp -a systemd/environ/rules %{buildroot}%{_sysconfdir}/%{shortname}/


install -d -m 0755 %{buildroot}%{_unitdir}
install -m 0644 -t %{buildroot}%{_unitdir} systemd/prometheus.service

# package infomation
%files
# set file attribute here
# need not list every file here, keep it as this
%{_prefix}

%{_bindir}/prometheus
%{_unitdir}/prometheus.service
%dir %{_sysconfdir}/%{shortname}
%config(noreplace) %{_sysconfdir}/%{shortname}/prometheus
%config(noreplace) %{_sysconfdir}/%{shortname}/prometheus.yml
%config(noreplace) %{_sysconfdir}/%{shortname}/rules

## indicate the dir for crontab

# %{_crondir}

%pre

%post
%systemd_post prometheus

%preun
%systemd_preun prometheus

%postun
%systemd_postun

if [ $1 -eq 0 ]
then
    # rpm -e , before delete files
    service prometheus stop
    %{_prefix}/prometheus  remove
fi

%changelog
* Tue Apr 18 2017 guoan.qga@alibaba-inc.com
- add spec of prometheus
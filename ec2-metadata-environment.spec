Name:           ec2-metadata-environment
Version:        1.0.0
Release:        1%{?dist}
Summary:        ec2-metadata-environment is a utility for accessing EC2 metadata as environment variables

Group:          System Environment/Daemons
License:        MIT
URL:            https://github.com/olly/ec2-metadata-enviroment
Source0:        update-ec2-metadata-environment
Source1:        %{name}.service
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

BuildRequires:  systemd-units
Requires:       systemd

%description
ec2-metadata-environment is a small utility which fetches EC2 metadata from the running instance, and stores it as enviroment variables in a file; which can be queried later.

%install
mkdir -p %{buildroot}/%{_bindir}
cp %{SOURCE0} %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_unitdir}
cp %{SOURCE1} %{buildroot}/%{_unitdir}/

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%clean
rm -rf %{buildroot}

%files
%attr(755, root, root) %{_bindir}/update-ec2-metadata-environment
%attr(644, root, root) %{_unitdir}/%{name}.service

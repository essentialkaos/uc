################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Simple utility for counting unique lines
Name:            uc
Version:         0.0.2
Release:         0%{?dist}
Group:           Applications/System
License:         Apache License, Version 2.0
URL:             https://kaos.sh/uc

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.11

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for counting unique lines.

################################################################################

%prep
%setup -q

%build
export GOPATH=$(pwd)
go build src/github.com/essentialkaos/%{name}/%{name}.go

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name} %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/%{name}

################################################################################

%changelog
* Fri Jan 31 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Added option -m/--max for defining maximum unique lines to process
- Added option -d/--dist for printing data distribution

* Thu Jan 30 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- Initial build

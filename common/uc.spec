################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Simple utility for counting unique lines
Name:            uc
Version:         1.0.1
Release:         0%{?dist}
Group:           Applications/System
License:         Apache License, Version 2.0
URL:             https://kaos.sh/uc

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.17

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Simple utility for counting unique lines.

################################################################################

%prep
%setup -q

%build
export GOPATH=$(pwd)
pushd src/github.com/essentialkaos/%{name}
  go build -mod vendor -o $GOPATH/%{name} %{name}.go
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -dm 755 %{buildroot}%{_mandir}/man1

install -pm 755 %{name} %{buildroot}%{_bindir}/

./%{name} --generate-man > %{buildroot}%{_mandir}/man1/%{name}.1

%clean
rm -rf %{buildroot}

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_mandir}/man1/%{name}.1.*
%{_bindir}/%{name}

################################################################################

%changelog
* Tue Mar 29 2022 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Thu Oct 22 2020 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added possibility to define -m/--max option as number with K and M
- Added man page generation

* Fri Jan 31 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.2-0
- Added option -m/--max for defining maximum unique lines to process
- Added option -d/--dist for printing data distribution

* Thu Jan 30 2020 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- Initial build
